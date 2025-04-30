import { IsString, IsOptional } from 'class-validator';

export class Address {
  @IsString()
  address_line1!: string;

  @IsOptional()
  @IsString()
  address_line2?: string;

  @IsString()
  city!: string;

  @IsString()
  state!: string;

  @IsString()
  zip!: string;
}
