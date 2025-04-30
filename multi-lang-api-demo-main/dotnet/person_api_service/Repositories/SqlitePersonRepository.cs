using Microsoft.EntityFrameworkCore;
using PersonApiService.Models;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace PersonApiService.Repositories
{
    public class SqlitePersonRepository : IPersonRepository
    {
        private readonly PersonContext _context;
        public SqlitePersonRepository(PersonContext context)
        {
            _context = context;
        }

        public async Task<IEnumerable<Person>> GetAllAsync() => await _context.Persons.ToListAsync();
        public async Task<Person?> GetByIdAsync(string id)
        {
            return await _context.Persons.FirstOrDefaultAsync(p => p.Id == id);
        }
        public async Task<IEnumerable<Person>> SearchByNameAsync(string? firstName, string? lastName)
        {
            var query = _context.Persons.AsQueryable();
            if (!string.IsNullOrEmpty(firstName))
                query = query.Where(p => EF.Functions.Like(p.FirstName, $"%{firstName}%"));
            if (!string.IsNullOrEmpty(lastName))
                query = query.Where(p => EF.Functions.Like(p.LastName, $"%{lastName}%"));
            return await query.ToListAsync();
        }
        public async Task<IEnumerable<Person>> ListByCityStateAsync(string city, string state) =>
            await _context.Persons.Where(p => p.Address.City == city && p.Address.State == state).ToListAsync();
        public async Task<Person> CreateAsync(Person person)
        {
            _context.Persons.Add(person);
            await _context.SaveChangesAsync();
            return person;
        }
        public async Task<bool> UpdateAsync(string id, Person person)
        {
            var existing = await _context.Persons.FirstOrDefaultAsync(p => p.Id == id);
            if (existing == null) return false;
            existing.FirstName = person.FirstName;
            existing.LastName = person.LastName;
            existing.Age = person.Age;
            existing.Address = person.Address;
            await _context.SaveChangesAsync();
            return true;
        }
        public async Task<bool> DeleteAsync(string id)
        {
            var person = await _context.Persons.FirstOrDefaultAsync(p => p.Id == id);
            if (person == null) return false;
            _context.Persons.Remove(person);
            await _context.SaveChangesAsync();
            return true;
        }
    }

    public class PersonContext : DbContext
    {
        public PersonContext(DbContextOptions<PersonContext> options) : base(options) { }
        public DbSet<Person> Persons { get; set; }
        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<Person>().ToTable("persons");
            modelBuilder.Entity<Person>().OwnsOne(p => p.Address, a =>
            {
                a.Property(ad => ad.AddressLine1).HasColumnName("address_line1");
                a.Property(ad => ad.AddressLine2).HasColumnName("address_line2");
                a.Property(ad => ad.City).HasColumnName("city");
                a.Property(ad => ad.State).HasColumnName("state");
                a.Property(ad => ad.Zip).HasColumnName("zip");
            });
        }
    }
}
