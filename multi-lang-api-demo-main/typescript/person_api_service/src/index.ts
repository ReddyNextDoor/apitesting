import { createApp } from './app-init';

const PORT = process.env.PORT || 8000;

(async () => {
  const app = await createApp();
  app.listen(PORT, () => {
    console.log(`Person API Service listening on port ${PORT}`);
  });
})();
