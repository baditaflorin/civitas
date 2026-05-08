# Privacy

Civitas ships with no client analytics.

The GitHub Pages frontend stores only non-sensitive UI preferences in browser `localStorage`, such as the configured backend endpoint. It does not store evidence files, extracted text, private entities, or secrets.

Evidence is sent only to the backend URL configured by the investigator. In production, that backend should be self-hosted, access-controlled at the network layer, and deployed with TLS.

Repository:

https://github.com/baditaflorin/civitas
