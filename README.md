# What is this?

This program let's you remove metadata components from installed unlocked packages.

# Why would I care?

If you are working with unlocked packages and want to move around metadata components between packages, or even get rid off unlocked packages completely, you will need to manually remove package components in the Salesforce UI (Installed Packages -> Select Package -> View components -> Remove).

This manual process can take forever if your package incoudes hundreds or even thousands of components.
But isn't there an official API or process to do exactly that you might wonder? Unfortunately not, at least not at this time of writing.

# How does it work then?

The program basically replicates the **Remove** clicks you would be doing in the UI to remove package components. However, in order to do so there are three pieces of information that you will need to provide (I know, and here you thought this would be easy..):

### My Domain URL

You can find your **My Domain URL** in the Setup under `Company Settings -> My Domain` (`/lightning/setup/OrgDomain/home`).

It will look similar to:

* Sandbox: `<company--envName>.my.salesforce.com`
* Scratch Org: `<randomName>.cs<regionNumber>.my.salesforce.com`
* Production: `<company>.my.salesforce.com`

Note: The URL can look a bit different in case **Enhanced Domains** are enabled.

### Valid Session ID

There are many ways to obtain a Session ID and unfortunately using the one from the SOAP Login will not work. In order to make sure you got the correct one (yes, there are different ones), I recommend to open one of the installed unlocked packages and click on **View Components**. Your URL should then look like this:

`https://<Domain>.lightning.force.com/lightning/setup/ImportedPackage/<PackageId>/Components/view`

Open up the Browser's Developer Tools and retrieve the `sid` cookie from the Domain that ends in `*.my.salesforce.com`. This is essential as there can be different session ids across domains.

### Unlocked Package Confirmation Token

Whenever you were to click manually on the **Remove** link to remove a package component, a confirmation token is passed along in the request. This token can be found by simply inspecting one of the **Remove** link elements in the Developer Console:

```html
<a href="javascript:if%20%28window.confirm%28%27Are%20you%20sure%20you%20want%20to%20remove%20this%20component%20from%20this%20package%3F%20Removing%20the%20component%20will%20not%20delete%20the%20component%20from%20your%20org.%20Please%20inform%20the%20owner%20of%20this%20package%20of%20this%20change%20so%20that%20they%20can%20make%20necessary%20changes%20to%20the%20package.%27%29%29%20window.location%3D%27%2F0A33O000000A8bj%3Fisdtp%3Dp1%26retURL%3D%252F0A33O000000A8bj%26p15%3D03d3O000000LfJz%26remove_package_member%3D1%26_CONFIRMATIONTOKEN%3DVmpFPSxNakF5TWkwd05DMHlOMVF5TURveE56b3pPQzQwTWpWYSwyUmMtWFoxbURxcENHTVFnVGJObmVpLE9UZGhNRGht%27%3B" class="actionLink" title="Remove - Record 1 - YourApexClassOrWhateverComponent">Remove</a>
```

As you can see the confirmation token hides in there at the end:

`VmpFPSxNakF5TWkwd05DMHlOMVF5TURveE56b3pPQzQwTWpWYSwyUmMtWFoxbURxcENHTVFnVGJObmVpLE9UZGhNRGht`

Note: While you only have to look for the Session ID once, you must retrieve the confirmation token for **every** package separately.

With these two things you are ready to mass remove unlocked packages components. Either pull the code and run it yourself (GO v1.18) or simply download the binary from the releases page.

# Usage

```bash
$ ./sfdxunpack --help
  -domain string
        My Domain URL
  -sid string
        Session ID
```

If no arguments are provided you will be prompted to provide the necessary information, otherwise you can provide everything when running the command:

```bash
$ ./sfdxunpack --domain=company.my.salesforce.com --sid='<SuperLongSessionIdThatShouldBeEnclosedInSingleQuotes>'

Unlocked packages: 
(1) Package A
(2) Package B
Select a package number: 1
Confirmation Token: <YourPackageConfirmationToken>

# Package components: 100
removing component '01p3O0000000000000' from package '0A33O0000000000000': 200 OK
99 remaining
```
