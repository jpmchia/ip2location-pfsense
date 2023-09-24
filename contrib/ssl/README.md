
See the following for more information:

https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/

https://stackoverflow.com/questions/44550970/firefox-54-stopped-trusting-self-signed-certs/48791236#48791236

To generate a self-signed certificate, run the following command:

openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]