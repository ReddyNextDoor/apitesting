import { SqlitePersonRepository } from './person.repository.sqlite';
import { MongoPersonRepository } from './person.repository.mongo';
import { IPersonRepository } from './person.repository.interface';

export function getPersonRepository(): IPersonRepository {
  const backend = (process.env.PERSON_REPO_BACKEND || 'sqlite').toLowerCase();
  if (backend === 'mongo') {
    return new MongoPersonRepository();
  }
  return new SqlitePersonRepository();
}
