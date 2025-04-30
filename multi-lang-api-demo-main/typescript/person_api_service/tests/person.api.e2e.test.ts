import request from 'supertest';
import { createApp } from '../src/app-init';

let app: any;
beforeAll(async () => {
  app = await createApp();
}, 30000); // 30s timeout for MongoDB startup

afterAll(async () => {
  // Close mongoose connection if used
  try {
    const mongoose = require('mongoose');
    if (mongoose.connection && mongoose.connection.readyState !== 0) {
      await mongoose.disconnect();
    }
  } catch (e) { /* ignore */ }
  // Close TypeORM connection if used
  try {
    const { AppDataSource } = require('../src/sqlite/data-source');
    if (AppDataSource && AppDataSource.isInitialized) {
      await AppDataSource.destroy();
    }
  } catch (e) { /* ignore */ }
});

describe('Person API E2E', () => {
  let createdId: string;
  const person = {
    first_name: 'John',
    last_name: 'Doe',
    age: 30,
    address: {
      address_line1: '123 Main St',
      city: 'Springfield',
      state: 'IL',
      zip: '62701',
    },
  };

  it('POST /persons creates a person', async () => {
    const res = await request(app).post('/persons').send(person);
    expect(res.status).toBe(201);
    expect(res.body).toHaveProperty('id');
    createdId = res.body.id || res.body._id;
  });

  it('GET /persons returns all persons', async () => {
    const res = await request(app).get('/persons');
    expect(res.status).toBe(200);
    expect(Array.isArray(res.body)).toBe(true);
  });

  it('GET /persons/:id returns the created person', async () => {
    const res = await request(app).get(`/persons/${createdId}`);
    expect(res.status).toBe(200);
    expect(res.body.first_name).toBe(person.first_name);
  });

  it('PUT /persons/:id updates the person', async () => {
    const updated = { ...person, age: 31 };
    const res = await request(app).put(`/persons/${createdId}`).send(updated);
    expect(res.status).toBe(200);
    const getRes = await request(app).get(`/persons/${createdId}`);
    expect(getRes.body.age).toBe(31);
  });

  it('GET /search?first_name=John returns the person', async () => {
    const res = await request(app).get('/search?first_name=John');
    expect(res.status).toBe(200);
    expect(res.body.some((p: any) => p.first_name === 'John')).toBe(true);
  });

  it('GET /citystate?city=Springfield&state=IL returns the person', async () => {
    const res = await request(app).get('/citystate?city=Springfield&state=IL');
    expect(res.status).toBe(200);
    expect(res.body.some((p: any) => p.address.city === 'Springfield')).toBe(true);
  });

  it('DELETE /persons/:id deletes the person', async () => {
    const res = await request(app).delete(`/persons/${createdId}`);
    expect(res.status).toBe(200);
    const getRes = await request(app).get(`/persons/${createdId}`);
    expect(getRes.status).toBe(404);
  });
});
