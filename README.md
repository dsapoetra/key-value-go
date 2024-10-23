# In-Memory Key-Value Store

A thread-safe in-memory key-value store implementation in Go, similar to Redis but with specific attribute-value pair storage capabilities.

## Features

- Thread-safe operations
- Type-safe attribute storage
- Support for string, numeric (float), and boolean values
- Command-line interface
- Sorted key retrieval
- Attribute-based search functionality

## Installation

1. Ensure you have Go installed on your system (version 1.16 or higher recommended)
2. Clone this repository:
```bash
git clone <repository-url>
cd key-value-store
```

3. Run the program:
```bash
go run main.go
```

## Commands

The application supports the following commands:

### PUT
Adds or updates a key with attribute-value pairs
```
put <key> <attributeKey1> <attributeValue1> <attributeKey2> <attributeValue2>...
```
Example:
```
put sde_bootcamp title SDE-Bootcamp price 30000.00 enrolled false estimated_time 30
```

### GET
Retrieves all attributes for a given key
```
get <key>
```
Example:
```
get sde_bootcamp
```
Output:
```
title: SDE-Bootcamp, price: 30000.00, enrolled: false, estimated_time: 30.0
```

### DELETE
Removes a key and its attributes from the store
```
delete <key>
```
Example:
```
delete sde_bootcamp
```

### SEARCH
Finds all keys that have a specific attribute key-value pair
```
search <attributeKey> <attributeValue>
```
Example:
```
search price 30000.00
```
Output:
```
sde_bootcamp
```

### KEYS
Lists all keys in the store in sorted order
```
keys
```
Output:
```
sde_bootcamp,sde_kickstart
```

### EXIT
Exits the program
```
exit
```

## Data Type Rules

- String values: Any text value
- Numeric values: All numbers are stored as float64 (e.g., "30000.00", "4000.00")
- Boolean values: Must be "true" or "false"
- Once an attribute's type is set, it cannot be changed
- Data type consistency is enforced across all entries

## Example Usage Session

```
# Add a bootcamp course
put sde_bootcamp title SDE-Bootcamp price 30000.00 enrolled false estimated_time 30

# Add another course
put sde_kickstart title SDE-Kickstart price 4000.00 enrolled true estimated_time 8

# Get course details
get sde_bootcamp
> title: SDE-Bootcamp, price: 30000.00, enrolled: false, estimated_time: 30.0

# List all keys
keys
> sde_bootcamp,sde_kickstart

# Search for courses by price
search price 30000.00
> sde_bootcamp

# Search for enrolled courses
search enrolled true
> sde_kickstart

# Delete a course
delete sde_bootcamp

# Verify deletion
get sde_bootcamp
> No entry found for sde_bootcamp
```

## Error Handling

1. Data Type Error
```
# Will result in error if trying to change attribute type
put sde_bootcamp title SDE-Bootcamp price true  # Error: price was previously float
> Data Type Error
```

2. Invalid Commands
- Unknown commands are ignored
- Malformed commands are ignored
- Missing arguments are handled gracefully

## Technical Details

- All operations are thread-safe using sync.RWMutex
- Keys are stored in sorted order
- Search results are returned in sorted order
- Numeric values are stored with consistent decimal precision
- Boolean values are case-sensitive ("true" or "false")
- String comparisons are case-sensitive

## Limitations

- Keys must be strings
- Attribute keys must be strings
- No persistence (in-memory only)
- No transaction support
- No TTL (Time To Live) support
- No nested objects support

## Best Practices

1. Consistent Attribute Naming
   - Use clear, descriptive attribute names
   - Maintain consistent naming conventions

2. Data Type Consistency
   - Plan attribute types before adding data
   - Use consistent numeric precision

3. Key Naming
   - Use meaningful, unique keys
   - Consider using prefixes for different types of data

## Contributing

Feel free to submit issues and enhancement requests!

## License

[MIT License](LICENSE)