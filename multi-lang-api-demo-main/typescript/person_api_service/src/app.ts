import express from 'express';
import swaggerUi from 'swagger-ui-express';
import * as fs from 'fs';
import * as dotenv from 'dotenv';
import { getPersonRepository } from './repositories/person.repository.factory';
import { AppDataSource } from './sqlite/data-source';
import { connectMongo } from './mongo/mongoose-connect';
import { PersonCreate, PersonUpdate } from './models/person.schema';
import { validateOrReject } from 'class-validator';
import { plainToInstance } from 'class-transformer';
import * as yaml from 'yaml';
import path from 'path';

// Load env vars
dotenv.config();

const app = express();
app.use(express.json());

// // Swagger/OpenAPI setup
// const openapiPath = path.join(__dirname, '../openapi.yaml');
// let openapiSpec: Record<string, any> | null = null;
// try {
//   openapiSpec = yaml.load(fs.readFileSync(openapiPath, 'utf-8')) as Record<string, any>;
//   console.log('Loaded OpenAPI path:', openapiPath);
//   console.log('Loaded OpenAPI title:', openapiSpec.info?.title);
// } catch (err) {
//   console.error('Failed to load OpenAPI YAML:', err);
// }
// // load the openapi.yaml file
// const swaggerDocument = yaml.load(fs.readFileSync(openapiPath, 'utf-8'));
// app.use('/docs', swaggerUi.serve, swaggerUi.setup(swaggerDocument));


const openapiPath = path.join(__dirname, '../openapi.yaml');
const fileContent = fs.readFileSync(openapiPath, 'utf8');
const swaggerDocument = yaml.parse(fileContent);
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocument));


// Root redirect
app.get('/', (req, res) => res.redirect('/api-docs'));

// Health endpoint
app.get('/health', (req, res) => {
  res.json({
    service: 'person-api-service',
    status: 'ok',
    mode: process.env.PERSON_REPO_BACKEND || 'sqlite',
    time: new Date().toISOString(),
    database: 'ok', // Replace with real DB check if needed
  });
});

// Repository initialization
let repo: ReturnType<typeof getPersonRepository>;
(async () => {
  const backend = (process.env.PERSON_REPO_BACKEND || 'sqlite').toLowerCase();
  if (backend === 'mongo') {
    const uri = process.env.MONGO_URI;
    if (!uri) throw new Error('MONGO_URI must be set for mongo backend');
    await connectMongo(uri);
  } else {
    await AppDataSource.initialize();
  }
  repo = getPersonRepository();
})();

// CRUD endpoints
app.get('/persons', async (req, res) => {
  res.json(await repo.getAll());
});

app.get('/persons/:id', async (req, res) => {
  const person = await repo.getById(req.params.id);
  if (!person) return res.status(404).json({ detail: 'Not found' });
  res.json(person);
});

app.post('/persons', async (req, res) => {
  try {
    const dto = plainToInstance(PersonCreate, req.body);
    await validateOrReject(dto);
    const created = await repo.create(dto);
    res.status(201).json(created);
  } catch (err) {
    res.status(400).json({ detail: err });
  }
});

app.put('/persons/:id', async (req, res) => {
  try {
    const dto = plainToInstance(PersonUpdate, req.body);
    await validateOrReject(dto);
    const ok = await repo.update(req.params.id, dto);
    if (!ok) return res.status(404).json({ detail: 'Not found' });
    res.json({ success: true });
  } catch (err) {
    res.status(400).json({ detail: err });
  }
});

app.delete('/persons/:id', async (req, res) => {
  const ok = await repo.delete(req.params.id);
  if (!ok) return res.status(404).json({ detail: 'Not found' });
  res.json({ success: true });
});

// Search endpoints
app.get('/search', async (req, res) => {
  const { first_name, last_name } = req.query;
  res.json(await repo.searchByName(first_name as string, last_name as string));
});

app.get('/citystate', async (req, res) => {
  const { city, state } = req.query;
  if (!city || !state) return res.status(400).json({ detail: 'city and state required' });
  res.json(await repo.listByCityState(city as string, state as string));
});

export default app;
