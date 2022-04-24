package salesforce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Salesforce interface {
	GetUnlockedPackages() ([]InstalledSubscriberPackage, error)
	GetMetadataComponents(subscriberPackageId string) ([]string, error)
	RemovePackageMetadata(package_id string, id string, confirmation_token string) error
}

type salesforce struct {
	Sandbox     bool
	MyDomainURL string
	SessionId   string
}

type QueryResult[T any] struct {
	Records []T `json:"records"`
}

type InstalledSubscriberPackage struct {
	Id                         string
	SubscriberPackageId        string
	SubscriberPackage          SubscriberPackage
	SubscriberPackageVersionId string
	SubscriberPackageVersion   SubscriberPackageVersion
}

type SubscriberPackage struct {
	Name string
}

type SubscriberPackageVersion struct {
	Package2ContainerOptions string
}

type Package2Member struct {
	Id                  string
	SubjectId           string
	SubscriberPackageId string
}

func New(myDomainURL, sessionID string, sandbox bool) Salesforce {
	if strings.Contains(myDomainURL, "https://") {
		myDomainURL = strings.TrimPrefix(myDomainURL, "https://")
	}
	return &salesforce{
		MyDomainURL: myDomainURL,
		SessionId:   sessionID,
		Sandbox:     sandbox,
	}
}

func (s *salesforce) GetUnlockedPackages() ([]InstalledSubscriberPackage, error) {
	client := &http.Client{}
	query := "SELECT Id, SubscriberPackageId, SubscriberPackage.Name, SubscriberPackageVersionId, SubscriberPackageVersion.IsManaged, SubscriberPackageVersion.Package2ContainerOptions FROM InstalledSubscriberPackage"
	reqURL := fmt.Sprintf("https://%s/services/data/v54.0/tooling/query/?q=%s", s.MyDomainURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %s", resp.Status)
	}

	var queryResults QueryResult[InstalledSubscriberPackage]
	if err := json.NewDecoder(resp.Body).Decode(&queryResults); err != nil {
		return nil, err
	}

	var unlockedPackages []InstalledSubscriberPackage
	for _, pkg := range queryResults.Records {
		if pkg.SubscriberPackageVersion.Package2ContainerOptions == "Unlocked" {
			unlockedPackages = append(unlockedPackages, pkg)
		}
	}
	return unlockedPackages, nil
}

func (s *salesforce) GetMetadataComponents(subscriberPackageId string) ([]string, error) {
	client := &http.Client{}
	query := fmt.Sprintf("SELECT Id, SubjectId, SubjectManageableState, SubscriberPackageId FROM Package2Member WHERE SubscriberPackageId = '%s'", subscriberPackageId)
	reqURL := fmt.Sprintf("https://%s/services/data/v54.0/tooling/query/?q=%s", s.MyDomainURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %s", resp.Status)
	}

	var queryResults QueryResult[Package2Member]
	if err := json.NewDecoder(resp.Body).Decode(&queryResults); err != nil {
		return nil, err
	}

	var metadataIds []string
	for _, r := range queryResults.Records {
		metadataIds = append(metadataIds, r.SubjectId)
	}
	return metadataIds, nil
}

func (s *salesforce) RemovePackageMetadata(package_id string, id string, confirmation_token string) error {
	fmt.Printf("removing component '%s' from package '%s'", id, package_id)

	client := &http.Client{}
	reqURL := fmt.Sprintf("https://%s/%s?isdtp=p1&p15=%s&remove_package_member=1&_CONFIRMATIONTOKEN=%s", s.MyDomainURL, package_id[:15], id[:15], confirmation_token)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Cookie", "sid="+s.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with status %s", resp.Status)
	}
	fmt.Printf(": %s\n", resp.Status)
	return nil
}
