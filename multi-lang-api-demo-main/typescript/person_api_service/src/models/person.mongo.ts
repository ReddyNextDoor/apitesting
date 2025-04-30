import { Schema, model, Document } from 'mongoose';

export interface Address {
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  zip: string;
}

export interface PersonDoc extends Document {
  first_name: string;
  last_name: string;
  age: number;
  address: Address;
}

const AddressSchema = new Schema<Address>({
  address_line1: { type: String, required: true },
  address_line2: { type: String },
  city: { type: String, required: true },
  state: { type: String, required: true },
  zip: { type: String, required: true },
}, { _id: false });

const PersonSchema = new Schema<PersonDoc>({
  first_name: { type: String, required: true },
  last_name: { type: String, required: true },
  age: { type: Number, required: true },
  address: { type: AddressSchema, required: true },
});

export const PersonModel = model<PersonDoc>('Person', PersonSchema);
