import { Entity, PrimaryGeneratedColumn, Column, OneToOne, JoinColumn } from 'typeorm';
import { Address } from './address.entity';

@Entity('persons')
export class Person {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  first_name!: string;

  @Column()
  last_name!: string;

  @Column()
  age!: number;

  @OneToOne(() => Address, { cascade: true, eager: true })
  @JoinColumn()
  address!: Address;
}
