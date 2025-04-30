import { AppDataSource } from '../sqlite/data-source';
import { Person } from '../models/person.entity';
import { Address } from '../models/address.entity';
import { IPersonRepository } from './person.repository.interface';
import { PersonCreate, PersonUpdate } from '../models/person.schema';
import { Like, Repository } from 'typeorm';

export class SqlitePersonRepository implements IPersonRepository {
  private personRepo: Repository<Person>;
  private addressRepo: Repository<Address>;

  constructor() {
    this.personRepo = AppDataSource.getRepository(Person);
    this.addressRepo = AppDataSource.getRepository(Address);
  }

  async getAll() {
    return await this.personRepo.find();
  }

  async getById(id: string) {
    return await this.personRepo.findOne({ where: { id } });
  }

  async searchByName(firstName?: string, lastName?: string) {
    const where: any = {};
    if (firstName) where.first_name = Like(`%${firstName}%`);
    if (lastName) where.last_name = Like(`%${lastName}%`);
    return await this.personRepo.find({ where });
  }

  async listByCityState(city: string, state: string) {
    return await this.personRepo
      .createQueryBuilder('person')
      .leftJoinAndSelect('person.address', 'address')
      .where('address.city = :city AND address.state = :state', { city, state })
      .getMany();
  }

  async create(person: PersonCreate) {
    const address = this.addressRepo.create(person.address);
    await this.addressRepo.save(address);
    const entity = this.personRepo.create({ ...person, address });
    return await this.personRepo.save(entity);
  }

  async update(id: string, person: PersonUpdate) {
    const entity = await this.personRepo.findOne({ where: { id } });
    if (!entity) return false;
    entity.first_name = person.first_name;
    entity.last_name = person.last_name;
    entity.age = person.age;
    if (entity.address) {
      Object.assign(entity.address, person.address);
      await this.addressRepo.save(entity.address);
    }
    await this.personRepo.save(entity);
    return true;
  }

  async delete(id: string) {
    const result = await this.personRepo.delete(id);
    return result.affected !== 0;
  }
}
