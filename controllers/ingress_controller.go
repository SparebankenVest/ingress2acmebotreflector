/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	ingress "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
	"time"
)

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type KeyVaultInfo struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	DNSNames          []string  `json:"dnsNames"`
	CreatedOn         time.Time `json:"createdOn"`
	ExpiresOn         time.Time `json:"expiresOn"`
	X509Thumbprint    string    `json:"x509Thumbprint"`
	KeyType           string    `json:"keyType"`
	KeySize           int       `json:"keySize"`
	KeyCurveName      *string   `json:"keyCurveName"`
	ReuseKey          bool      `json:"reuseKey"`
	IsExpired         bool      `json:"isExpired"`
	IsIssuedByAcmebot bool      `json:"isIssuedByAcmebot"`
	IsSameEndpoint    bool      `json:"isSameEndpoint"`
	AcmeEndpoint      string    `json:"acmeEndpoint"`
}

type CertOrder struct {
	DnsNames []string `json:"DnsNames"`
}



var api_scope = os.Getenv("API_SCOPE")
var backend = os.Getenv("BACKEND")
var acmebot_rest_api_timeout = os.Getenv("ACMEBOT_REST_API_TIMEOUT") // default
var azure_ad_client_id = os.Getenv("AZURE_AD_CLIENT_ID")
var domains = strings.Split(os.Getenv("DOMAINS"), ",")
var timeout int = 20 // default

func init() {
	var err error 
	 timeout, err = strconv.Atoi(acmebot_rest_api_timeout);
	if err != nil {
		// handle error, e.g., log it or set a default value
		logger.Error(err, "Parsing ACMEBOT_REST_API_TIMEOUT failed, using default value")

	}
}

var myClient = &http.Client{Timeout: time.Duration(timeout) * time.Second}
var logger = log.Log.WithName("ingress_controller")
var s = flag.String("s", api_scope+"/.default", "scope for access token")

func getJson(url string, target interface{}, token string) error {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}
	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func checkExistingCerts(dnsName string, keyVaultInfoList []KeyVaultInfo) bool {
	logger.Info("Found " + strconv.Itoa(len(keyVaultInfoList)) + " certs in the keyvault")
	for _, keyVaultInfo := range keyVaultInfoList {
		for _, dns := range keyVaultInfo.DNSNames {
			if dns == dnsName {
				return true
			}
		}
	}
	return false
}

//+kubebuilder:rbac:groups=core,resources=ingresses/finalizers,verbs=get;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)
	var ing ingress.Ingress

	flag.Parse()
	opts := &azidentity.ManagedIdentityCredentialOptions{}
	opts.ID = azidentity.ClientID(azure_ad_client_id)

	cred, err := azidentity.NewManagedIdentityCredential(opts)
	if err != nil {
		logger.Error(err, "unable to create managed identity credential")
		return ctrl.Result{}, nil
	}

	logger.Info("Calling GetToken()...")
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{*s}})
	if err != nil {
		logger.Error(err, "unable to get token")
		return ctrl.Result{}, nil
	}

	logger.Info("Reconciling Ingress")
	if err := r.Get(ctx, req.NamespacedName, &ing); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Ingress resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if ing.Spec.TLS == nil {
		logger.Info("No TLS found. Ignoring")
		return ctrl.Result{}, nil
	}
	secretName := ing.Spec.TLS[0].SecretName
	host := ing.Spec.Rules[0].Host
	logger.Info("TLS found, secret name: " + secretName + " host: " + host)

	var keyVaultInfoList []KeyVaultInfo

	//check if the current host is in the list of supported domains
	var validDomain bool
	for _, domain := range domains {
		if strings.Contains(host, domain) {
			validDomain = true
		}
	}
	if !validDomain {
		logger.Info("Host not in list of supported domains, exiting")
		return ctrl.Result{}, nil
	}

	err = getJson(backend+"/api/certificates", &keyVaultInfoList, tk.Token)
	if err != nil {
		logger.Error(err, "unable to get cert list")
		return ctrl.Result{}, err
	}

	logger.Info("checking existing certs")
	existingCerts := checkExistingCerts(host, keyVaultInfoList)

	//if no existing certs, order new cert through the rest api
	if !existingCerts {
		logger.Info("No cert found in keyvault, initiating order")
		certorder := CertOrder{DnsNames: []string{host}}

		payloadBuf := new(bytes.Buffer)

		json.NewEncoder(payloadBuf).Encode(certorder)

		req, err := http.NewRequest("POST", backend+"/api/certificate", payloadBuf)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tk.Token)

		response, err := myClient.Do(req)
		if err != nil {
			logger.Error(err, "unable to order new cert")
			return ctrl.Result{}, err
		}
		defer response.Body.Close()
		println("Status: ", response.Status)
		certorderUrl := response.Header.Get("Location")
		logger.Info("Cert ordered, waiting for 30 seconds, callback: " + certorderUrl)
		time.Sleep(30 * time.Second)
		logger.Info("Checking if cert is ready")
		req2, err := http.NewRequest("GET", certorderUrl, nil)
		resp, err := myClient.Do(req2)
		req2.Header.Set("Authorization", "Bearer "+tk.Token)

		if err != nil {
			logger.Error(err, "unable to get cert status")
			return ctrl.Result{}, err
		} else if resp.StatusCode == 200 {
			logger.Info("Cert is ready")
		}
		defer resp.Body.Close()
	} else {
		logger.Info("Cert found in keyvault, exiting")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ingress.Ingress{}).
		Complete(r)
}
