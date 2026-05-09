# Privacy

Civitas ships with no client analytics.

The GitHub Pages frontend stores only local UI session preferences in browser `localStorage`:

- Configured backend endpoint.
- Selected case id.
- Search term.

The frontend does not store evidence files, extracted text, entities, private notes, or secrets. Evidence is sent only to the backend URL configured by the investigator.

Downloaded safe exports and downloaded case state files can contain sensitive evidence metadata or source bytes. Treat those files like the original evidence dump, store them encrypted where appropriate, and do not attach them to public issues.

In production, the backend should be self-hosted, access-controlled at the network layer, and deployed with TLS.

Repository:

https://github.com/baditaflorin/civitas
