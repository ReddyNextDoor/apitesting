import { Schema, model, Document } from 'mongoose';

export interface AddressDoc extends Document {
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  zip: string;
}

const AddressSchema = new Schema<AddressDoc>({
  address_line1: { type: String, required: true },
  address_line2: { type: String },
  city: { type: String, required: true },
  state: { type: String, required: true },
  zip: { type: String, required: true },
});

export const AddressModel = model<AddressDoc>('Address', AddressSchema);
