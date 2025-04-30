import { PersonModel, PersonDoc } from '../models/person.mongo';
import { IPersonRepository } from './person.repository.interface';
import { PersonCreate, PersonUpdate } from '../models/person.schema';
import mongoose from 'mongoose';

export class MongoPersonRepository implements IPersonRepository {
  async getAll() {
    return await PersonModel.find().exec();
  }

  async getById(id: string) {
    return await PersonModel.findById(id).exec();
  }

  async searchByName(firstName?: string, lastName?: string) {
    const filter: any = {};
    if (firstName) filter.first_name = new RegExp(firstName, 'i');
    if (lastName) filter.last_name = new RegExp(lastName, 'i');
    return await PersonModel.find(filter).exec();
  }

  async listByCityState(city: string, state: string) {
    // Find persons whose embedded address matches city and state
    return await PersonModel.find({
      'address.city': city,
      'address.state': state
    }).exec();
  }

  async create(person: PersonCreate) {
    const personDoc = new PersonModel({ ...person, address: person.address });
    await personDoc.save();
    return personDoc;
  }

  async update(id: string, person: PersonUpdate) {
    const existing = await PersonModel.findById(id).exec();
    if (!existing) return false;
    await PersonModel.findByIdAndUpdate(id, {
      first_name: person.first_name,
      last_name: person.last_name,
      age: person.age,
      address: person.address
    }, { new: true });
    return true;
  }

  async delete(id: string) {
    const existing = await PersonModel.findById(id).exec();
    if (!existing) return false;
    await PersonModel.findByIdAndDelete(id);
    return true;
  }
}
