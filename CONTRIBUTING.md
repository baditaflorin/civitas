# Contributing

Thanks for helping improve Civitas.

## Local Setup

```sh
make install-hooks
make dev
```

Before pushing, run:

```sh
make fmt
make lint
make test
make build
make smoke
```

Use Conventional Commits:

- `feat: add evidence timeline`
- `fix: handle empty upload batches`
- `docs: update deployment guide`
- `data: refresh sample corpus`

Never commit secrets, private evidence, personal data, or live `.env` files.
