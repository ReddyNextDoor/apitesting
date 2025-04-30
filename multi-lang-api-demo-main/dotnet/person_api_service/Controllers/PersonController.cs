using Microsoft.AspNetCore.Mvc;
using PersonApiService.Models;
using PersonApiService.Repositories;

namespace PersonApiService.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class PersonController : ControllerBase
    {
        private readonly IPersonRepository _repo;
        public PersonController(IPersonRepository repo)
        {
            _repo = repo;
        }

        [HttpGet]
        public async Task<IActionResult> GetAll() => Ok(await _repo.GetAllAsync());

        [HttpGet("{id}")]
        public async Task<IActionResult> GetById(string id)
        {
            var person = await _repo.GetByIdAsync(id);
            return person is null ? NotFound() : Ok(person);
        }

        [HttpGet("search")] // /api/person/search?firstName=...&lastName=...
        public async Task<IActionResult> SearchByName([FromQuery] string? firstName = null, [FromQuery] string? lastName = null)
            => Ok(await _repo.SearchByNameAsync(firstName, lastName));

        [HttpGet("by_city_state")]
        public async Task<IActionResult> ListByCityState([FromQuery] string city, [FromQuery] string state)
            => Ok(await _repo.ListByCityStateAsync(city, state));

        [HttpPost]
        public async Task<IActionResult> Create([FromBody] Person person)
        {
            // If Id is null or empty, set it to a new GUID (for SQLite); MongoDB will generate its own if not set
            if (string.IsNullOrEmpty(person.Id))
                person.Id = Guid.NewGuid().ToString();
            var created = await _repo.CreateAsync(person);
            return CreatedAtAction(nameof(GetById), new { id = created.Id }, created);
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> Update(string id, [FromBody] Person person)
        {
            var updated = await _repo.UpdateAsync(id, person);
            return updated ? NoContent() : NotFound();
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> Delete(string id)
        {
            var deleted = await _repo.DeleteAsync(id);
            return deleted ? NoContent() : NotFound();
        }
    }
}
