import { Entity, PrimaryGeneratedColumn, Column, OneToOne } from 'typeorm';
import { Address as AddressDTO } from './address.schema';

@Entity('addresses')
export class Address {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  address_line1!: string;

  @Column({ nullable: true })
  address_line2?: string;

  @Column()
  city!: string;

  @Column()
  state!: string;

  @Column()
  zip!: string;
}
