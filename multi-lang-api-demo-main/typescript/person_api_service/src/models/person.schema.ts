import { IsString, IsInt, ValidateNested } from 'class-validator';
import { Type } from 'class-transformer';
import { Address } from './address.schema';

export class PersonBase {
  @IsString()
  first_name!: string;

  @IsString()
  last_name!: string;

  @IsInt()
  age!: number;

  @ValidateNested()
  @Type(() => Address)
  address!: Address;
}

export class PersonCreate extends PersonBase {}
export class PersonUpdate extends PersonBase {}

export class PersonOut extends PersonBase {
  id!: string;
}
