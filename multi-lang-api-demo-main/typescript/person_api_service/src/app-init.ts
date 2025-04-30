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

dotenv.config();

function normalizePerson(obj: any) {
  if (!obj) return obj;
  const plain = obj.toObject ? obj.toObject() : obj;
  if (plain._id && !plain.id) plain.id = plain._id.toString();
  delete plain._id;
  delete plain.__v;
  if (plain.address && plain.address._id) delete plain.address._id;
  return plain;
}

export async function createApp() {
  const app = express();
  app.use(express.json());

  // Swagger/OpenAPI setup
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
  const backend = (process.env.PERSON_REPO_BACKEND || 'sqlite').toLowerCase();
  if (backend === 'mongo') {
    const uri = process.env.MONGO_URI;
    if (!uri) throw new Error('MONGO_URI must be set for mongo backend');
    await connectMongo(uri);
  } else {
    await AppDataSource.initialize();
  }
  repo = getPersonRepository();

  // CRUD endpoints
  app.get('/persons', async (req, res) => {
    const people = await repo.getAll();
    res.json(Array.isArray(people) ? people.map(normalizePerson) : []);
  });

  app.get('/persons/:id', async (req, res) => {
    const person = await repo.getById(req.params.id);
    if (!person) return res.status(404).json({ detail: 'Not found' });
    res.json(normalizePerson(person));
  });

  app.post('/persons', async (req, res) => {
    try {
      const dto = plainToInstance(PersonCreate, req.body);
      await validateOrReject(dto);
      const created = await repo.create(dto);
      res.status(201).json(normalizePerson(created));
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
      // Return the updated person for consistency
      const updated = await repo.getById(req.params.id);
      res.json(normalizePerson(updated));
    } catch (err) {
      res.status(400).json({ detail: err });
    }
  });

  app.delete('/persons/:id', async (req, res) => {
    const ok = await repo.delete(req.params.id);
    if (!ok) return res.status(404).json({ detail: 'Not found' });
    res.json({ success: true });
  });

  app.get('/search', async (req, res) => {
    const { first_name, last_name } = req.query;
    const results = await repo.searchByName(first_name as string | undefined, last_name as string | undefined);
    res.json(Array.isArray(results) ? results.map(normalizePerson) : []);
  });

  app.get('/citystate', async (req, res) => {
    const { city, state } = req.query;
    if (!city || !state) return res.status(400).json({ detail: 'city and state are required' });
    const results = await repo.listByCityState(city as string, state as string);
    res.json(Array.isArray(results) ? results.map(normalizePerson) : []);
  });

  return app;
}
