# Examples with Hippo SDK

## Install dependencies

```bash
npm install
# Install wasm-pack built npm package
npm install ../../pkg/hippo-sdk-1.0.0.tgz
```

## Build latest SDK(optional)

```bash
# In repository Root
wasm-pack build --target=web
```

Move output `pkg` directory into `sdk/core/examples/web/src/assets`

## Run example

```bash
npm run dev
```
