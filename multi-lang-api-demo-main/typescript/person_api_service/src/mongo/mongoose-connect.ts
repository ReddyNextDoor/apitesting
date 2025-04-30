import mongoose from 'mongoose';

export async function connectMongo(uri: string) {
  await mongoose.connect(uri, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  } as any);
}
