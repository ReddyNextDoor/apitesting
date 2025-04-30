import { Person } from '../models/person.entity';
import { PersonDoc } from '../models/person.mongo';
import { PersonCreate, PersonUpdate } from '../models/person.schema';

export interface IPersonRepository {
  getAll(): Promise<any[]>;
  getById(id: string): Promise<any | null>;
  searchByName(firstName?: string, lastName?: string): Promise<any[]>;
  listByCityState(city: string, state: string): Promise<any[]>;
  create(person: PersonCreate): Promise<any>;
  update(id: string, person: PersonUpdate): Promise<boolean>;
  delete(id: string): Promise<boolean>;
}
