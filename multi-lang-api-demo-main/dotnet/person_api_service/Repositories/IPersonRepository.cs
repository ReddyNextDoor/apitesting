using System.Collections.Generic;
using System.Threading.Tasks;
using PersonApiService.Models;

namespace PersonApiService.Repositories
{
    public interface IPersonRepository
    {
        Task<IEnumerable<Person>> GetAllAsync();
        Task<Person?> GetByIdAsync(string id);
        Task<IEnumerable<Person>> SearchByNameAsync(string? firstName, string? lastName);
        Task<IEnumerable<Person>> ListByCityStateAsync(string city, string state);
        Task<Person> CreateAsync(Person person);
        Task<bool> UpdateAsync(string id, Person person);
        Task<bool> DeleteAsync(string id);
    }
}
