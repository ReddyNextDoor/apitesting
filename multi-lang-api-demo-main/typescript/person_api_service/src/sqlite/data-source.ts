import 'reflect-metadata';
import { DataSource } from 'typeorm';
import { Person } from '../models/person.entity';
import { Address } from '../models/address.entity';
import * as path from 'path';

const dbPath = process.env.SQLITE_DB_PATH || path.join(__dirname, '../../db/persons.db');

export const AppDataSource = new DataSource({
  type: 'sqlite',
  database: dbPath,
  entities: [Person, Address],
  synchronize: true,
});
