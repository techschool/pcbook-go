In the [previous article](https://dev.to/techschoolguru/a-complete-overview-of-ssl-tls-and-its-cryptographic-system-36pd), we’ve talked about how digital certificates help with authentication and provide a safe and reliable key exchange process in TLS. Today we will learn exactly how to generate a certificate and have it signed by a Certificate Authority (CA).

{% youtube sWLf2tyY2q0 %}

For the purpose of this demo, we won’t submit our Certificate Signing Request (CSR) to a real CA. Instead, we will play both roles: the certificate authority and the certificate applicant.

So in the first step, we will generate a private key and its self-signed certificate for the CA. They will be used to sign the CSR later. In the second step, we will generate a private key and its paired CSR for the web server that we want to use TLS. Then finally we will use the CA’s private key to sign the web server’s CSR and get back the signed certificate.

In order to do all of these things, We need to have openssl installed. If you’re on a Mac, it’s probably already there. You can run `openssl version` to see which version it’s running.

In my case, it’s [LibreSSL](https://www.libressl.org/) version 2.8.3. We can access its manual at [this link](https://man.openbsd.org/openssl.1)

The first command we’re gonna used is `openssl req`, which stands for request.

This command is used to create and process certificate signing request. It can also be used to create a self-signed certificate for the CA, which is exactly what we want in the first step.

```bash
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem
```

The `-x509` option is used to tell openssl to output a self-signed certificate instead of a certificate request. In case you don’t know, [X509](https://en.wikipedia.org/wiki/X.509) is just a standard format of the public key certificate. 

The `-newkey rsa:4096` option basically tells openssl to create both a new RSA private key (4096-bit) and its certificate request at the same time. As we’re using this together with `-x509` option, it will output a certificate instead of a cetificate request.

The next option is `-days 365`, which specifies the number of days that the certificate is valid for.

Then we use the `-keyout` option to tell openssl to write the created private key to `ca-key.pem` file

And finally the `-out` option to tell it to write the certificate to `ca-cert.pem` file.

When we run this command, openssl will start generating the private key. 

Once the key is generated, we will be asked to provide a pass phrase, which will be used to encrypt the private key before writing it to the PEM file.

Why is it encrypted? Because if somehow the private key file is hacked, the hacker cannot use it to do anything without knowing the pass phrase to decrypt it first.

Next, openssl will ask us for some identity information to generate the certificate:

- The country code
- The state or province name
- The organization name
- The unit name
- The common name (or domain name)
- The email address.

And that’s it! The certificate and private key files will be successfully generated.

If we cat the private key file, we can see it says `ENCRYPTED PRIVATE KEY`.

The certificate, on the other hand, is not encrypted, but only base64-encoded, because it just contains the public key, the identity information, and the signature that should be visible to everyone.

We can use the `openssl x509` command to display all the information encoded in this certificate. This command can also be used to sign certificate requests, which we will do in a few minute.

```bash
openssl x509 -in ca-cert.pem -noout -text
```

We use the `-in` option to pass in the CA’s certificate file.

And the `-noout` option to tell it to not output the original base64-encoded value.

Instead, we use the `-text` option because we want to display it in a readable text format.

Here we can see all information of the certificate, such as the version, the serial number. The issuer and the subject are the same in this case because this is a self-signed certificate. Then the RSA public key and signature.

I’m gonna copy this command and save it to a `gen.sh` script.

With this script, I want to automate the process of generating a set of keys and certificates.

Before moving to the 2nd step, I’m gonna show you another way to provide the identity information without entering it interactively as before.

To do this, we must add the `-subj` (subject) option to the `openssl req` command:

```bash
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=Tech School/OU=Education/CN=*.techschool.guru/emailAddress=techschool.guru@gmail.com"
```

Let’s add a command to remove all pem files at the top of the `gen.sh` script, and run it in the terminal.

We still being prompted for a pass phrase, but it doesn’t ask for identity information anymore, because we already provided them in the subject option. Great!

Now the next step is to generate a private key and CSR for our web server. 

It’s almost the same as the command we used in the 1st step. Except that, this time we don’t want to self-sign it, so we should remove this `-x509` option. The `-days` option should be removed as well, since we don’t create a certificate, but just a CSR.

```bash
openssl req -newkey rsa:4096 -keyout server-key.pem -out server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"
```

The name of the output key should be `server-key.pem`. The output certificate request file should be `server-req.pem`. And the subject should contain our web server’s information.

Now, when we run this command, the encrypted private key and the certificate signing request files will be generated.

This time, in the `server-req.pem` file, it says `CERTIFICATE REQUEST`, not `CERTIFICATE` as in the `ca-cert.pem` file. That's because it’s not a certificate as before, but a certificate signing request instead.

So now let’s move to step 3 and sign this request. For that, we will use the same `openssl x509` command that we’ve used to display certificate before.

Let’s open the terminal and run this:

```bash
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
```

In this command, we use the `-req` option to tell openssl that we’re gonna pass in a certificate request.

We use the `-in` option follow by the name of the request file: `server-req.pem`.

Next we use the `-CA` option to pass in the certificate file of the CA: `ca-cert.pem`.

And the `-CAkey` option to pass in the private key of the CA: `ca-key.pem`.

Then 1 important option is `-CAcreateserial`. Basically the CA must ensure that each certificate it signs goes with a unique serial number. So with this option, a file containing the next serial number will be generated if it doesn’t exist.

Finally we use the `-out` option to specify the file to write the output certificate to.

Now as you can see here, because the CA’s private key is encrypted, openssl is asking for the pass phrase to decrypt it before it can be used to sign the certificate. It’s a countermeasure in case the CA’s private key is hacked.

OK, now we’ve got the signed certificate for our web server. Let’s print it out in plain text format.

This is its unique serial number. And we can also see a `ca-cert.srl` file, which contains the same serial number.

The issuer section contains the information of the CA, which is Tech School in this case.

By default, the certificate is valid for 30 days. We can change it by adding the `-days` option to the signing command.

```bash
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
```

Now the validity duration has changed to 60 days.

If you remember the Youtube certificate that we’ve seen in the previous video, it is used for many Google websites with different domain names.

We can also do that for our web server by specifying the Subject Alternative Name extension when signing the certificate request.

The `-extfile` option allows us to state the file containing the extensions. We can see the format of the config file in [this page](https://man.openbsd.org/x509v3.cnf.5#Subject_alternative_name).

There are several things that we can use as the alternative name, such as email, DNS, or IP.

I will create a new file `server-ext.cnf` with this content:

```config
subjectAltName=DNS:*.pcbook.com,DNS:*.pcbook.org,IP:0.0.0.0
```

Here I set `DNS` to multiple domain names: `*.pcbook.com` and `*.pcbook.org`.

I also set `IP` to `0.0.0.0` which will be used when we develop on localhost.

Now in the certificate signing command, let’s add the -extfile option and pass in the name of the extension config file:

```bash
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
```

Now the result certificate file has a new extensions section with all the subject alternative names that we’ve chosen. Awesome!

So looks like our automate script is ready, except for the fact that we have to enter a lot of password to protect the private keys.

In case we just want to use this for development and testing, we can tell openssl to not encrypt the private key, so that it won’t ask us for the pass phrase.

We do that by adding the `-nodes` option to the `openssl req` command like this:

```bash
rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=Tech School/OU=Education/CN=*.techschool.guru/emailAddress=techschool.guru@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text
```

Now if we run `gen.sh` again, it will not ask for passwords anymore. 

And if we look at the private key files, it is now `PRIVATE KEY`, and not `ENCRYPTED PRIVATE KEY` as before. Cool!

One last thing before we finish, I will show you how to verify if a certificate is valid or not.

We can do that with the `openssl verify` command:

```bash
openssl verify -CAfile ca-cert.pem server-cert.pem
```

Pass in the trusted CA’s certificate and the certificate that we want to verify. If it returns `OK` then the certificate is valid.

And that’s it for today’s article. I hope it’s useful for you. Thanks for reading and I’ll see you guys in the next one!


