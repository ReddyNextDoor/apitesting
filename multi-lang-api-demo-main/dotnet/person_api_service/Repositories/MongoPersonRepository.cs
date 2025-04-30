using MongoDB.Driver;
using PersonApiService.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace PersonApiService.Repositories
{
    public class MongoPersonRepository : IPersonRepository
    {
        private readonly IMongoCollection<Person> _collection;
        public MongoPersonRepository(IMongoDatabase database)
        {
            _collection = database.GetCollection<Person>("persons");
        }
        public async Task<IEnumerable<Person>> GetAllAsync() =>
            await _collection.Find(_ => true).ToListAsync();
        public async Task<Person?> GetByIdAsync(string id) =>
            await _collection.Find(p => p.Id == id).FirstOrDefaultAsync();
        public async Task<IEnumerable<Person>> SearchByNameAsync(string? firstName, string? lastName)
        {
            var filter = Builders<Person>.Filter.Empty;
            if (!string.IsNullOrEmpty(firstName))
                filter &= Builders<Person>.Filter.Regex(p => p.FirstName, new MongoDB.Bson.BsonRegularExpression(firstName, "i"));
            if (!string.IsNullOrEmpty(lastName))
                filter &= Builders<Person>.Filter.Regex(p => p.LastName, new MongoDB.Bson.BsonRegularExpression(lastName, "i"));
            return await _collection.Find(filter).ToListAsync();
        }
        public async Task<IEnumerable<Person>> ListByCityStateAsync(string city, string state) =>
            await _collection.Find(p => p.Address.City == city && p.Address.State == state).ToListAsync();
        public async Task<Person> CreateAsync(Person person)
        {
            await _collection.InsertOneAsync(person);
            return person;
        }
        public async Task<bool> UpdateAsync(string id, Person person)
        {
            // Ensure the replacement document has the correct Id (MongoDB _id)
            person.Id = id;
            var result = await _collection.ReplaceOneAsync(p => p.Id == id, person);
            return result.ModifiedCount > 0;
        }
        public async Task<bool> DeleteAsync(string id)
        {
            var result = await _collection.DeleteOneAsync(p => p.Id == id);
            return result.DeletedCount > 0;
        }
    }
}
