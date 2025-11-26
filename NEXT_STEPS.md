# Next Steps - Week 2 Planning

## Overview
Week 1 has successfully established the foundation with user authentication, authorization, and microservices architecture. Week 2 will focus on implementing the core trip management and search functionality that enables users to discover and book bus trips.

---

## Week 2 Objectives

### Primary Goals
1. **Trip Search & Discovery** - Enable users to search for available trips with advanced filtering
2. **Route Management** - Allow admins to configure routes with multiple pickup/dropoff points
3. **Seat Configuration** - Implement visual seat map management with pricing tiers
4. **Performance Optimization** - Add caching and database indexes for search performance

---

## Planned Features

### 1. Route Configuration with Pickup/Dropoff Points

**Problem:** Current route model only supports single origin and destination. Real-world bus routes have multiple stops.

**Solution:**
- Create `RouteStop` model with GPS coordinates
- Support pickup-only, dropoff-only, or both types
- Time offset calculation from departure
- Reorderable stop sequences

**Technical Implementation:**
```
- New table: route_stops
- Fields: location, latitude, longitude, offset_minutes, stop_type
- API endpoints for CRUD operations
```

**Estimated Time:** 2 hours

---

### 2. Seat Map Configuration Tool

**Problem:** Current bus model only has total seat capacity. Need visual layout for seat selection.

**Solution:**
- Create `Seat` model with row/column positioning
- Support different seat types (Standard, VIP, Sleeper)
- Pricing multipliers per seat type
- Double-decker bus support (floor 1/2)

**Technical Implementation:**
```
- Enhanced seats table with layout fields
- Bulk seat creation from templates
- Visual seat map API response
- Price calculation: base_price × seat_multiplier
```

**Estimated Time:** 4 hours

---

### 3. Advanced Trip Search API

**Current State:** Basic search by origin, destination, and date.

**Planned Enhancements:**
- **Time Range Filter** - Search within specific departure hours
- **Price Range Filter** - Min/max price filtering
- **Amenities Filter** - WiFi, AC, Toilet, etc.
- **Seat Type Filter** - Filter by available seat types
- **Operator Filter** - Search by specific bus operators
- **Sorting** - By price, departure time, or duration
- **Pagination** - Efficient handling of large result sets

**Technical Implementation:**
```
- Update TripRepository with dynamic query building
- Add filter aggregations for UI
- Implement sorting logic
- Response includes: trips, filters, pagination metadata
```

**Estimated Time:** 3 hours

---

### 4. Type-Safe Constants

**Problem:** Magic strings throughout codebase ("standard", "vip", etc.) prone to typos.

**Solution:**
- Create constants package with typed enums
- Validation methods (IsValid())
- Helper methods (GetPriceMultiplier())

**Constants to Define:**
```go
- SeatType: Standard, VIP, Sleeper
- BusType: Standard, VIP, Sleeper, DoubleDecker
- TripStatus: Scheduled, InProgress, Completed, Cancelled
- StopType: Pickup, Dropoff, Both
- Amenity: WiFi, AC, Toilet, TV, Charging, etc.
```

**Estimated Time:** 1 hour

---

### 5. Performance Optimization

#### Database Indexes
```sql
- idx_trips_route_date: Optimize trip search by route and date
- idx_routes_origin_dest: Fast route lookup
- idx_trips_price: Price range filtering
- idx_route_stops_route: Stop lookup by route
```

**Estimated Time:** 1 hour

#### Redis Caching Strategy
```
- Search results: 5 min TTL
- Route details: 1 hour TTL
- Seat availability: 1 min TTL
- Cache invalidation on updates
```

**Estimated Time:** 2 hours

---

## File Organization

### New Files to Create
```
backend/trip-service/
├── internal/
│   ├── constants/
│   │   ├── seat_type.go
│   │   ├── bus_type.go
│   │   ├── trip_status.go
│   │   ├── stop_type.go
│   │   └── amenity.go
│   ├── model/
│   │   ├── route_stop.go
│   │   ├── seat.go (enhanced)
│   │   └── filters.go
│   ├── repository/
│   │   ├── route_stop_repository.go
│   │   └── seat_repository.go (enhanced)
│   ├── service/
│   │   ├── route_stop_service.go
│   │   ├── seat_service.go
│   │   └── cache_service.go
│   └── handler/
│       ├── route_stop_handler.go
│       └── seat_handler.go
└── migrations/
    ├── 001_create_route_stops.sql
    ├── 002_update_seats_table.sql
    └── 003_add_search_indexes.sql
```

---

## API Endpoints to Implement

### Route Stops
```
POST   /api/v1/routes/stops              - Create stop
PUT    /api/v1/routes/stops/:id          - Update stop
DELETE /api/v1/routes/stops/:id          - Delete stop
GET    /api/v1/routes/:route_id/stops    - List stops
POST   /api/v1/routes/:route_id/stops/reorder - Reorder stops
```

### Seat Management
```
POST   /api/v1/buses/seats               - Create seat
POST   /api/v1/buses/seats/bulk          - Bulk create
PUT    /api/v1/buses/seats/:id           - Update seat
DELETE /api/v1/buses/seats/:id           - Delete seat
GET    /api/v1/buses/:bus_id/seats       - Get seat map
```

### Enhanced Search
```
POST   /api/v1/trips/search              - Advanced search with filters
```

---

## Database Schema Changes

### New Table: route_stops
```sql
CREATE TABLE route_stops (
    id UUID PRIMARY KEY,
    route_id UUID REFERENCES routes(id),
    stop_order INTEGER NOT NULL,
    stop_type VARCHAR(20),  -- 'pickup', 'dropoff', 'both'
    location VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    offset_minutes INTEGER,
    is_active BOOLEAN DEFAULT true
);
```

### Enhanced Table: seats
```sql
ALTER TABLE seats ADD COLUMN row INTEGER;
ALTER TABLE seats ADD COLUMN "column" INTEGER;
ALTER TABLE seats ADD COLUMN seat_type VARCHAR(20);
ALTER TABLE seats ADD COLUMN price_multiplier DECIMAL(3,2) DEFAULT 1.0;
ALTER TABLE seats ADD COLUMN floor INTEGER DEFAULT 1;
```

---

## Testing Strategy

### Unit Tests
- Route stop service CRUD operations
- Seat service with pricing calculations
- Cache service hit/miss scenarios
- Search filters validation

### Integration Tests
- End-to-end trip search flow
- Seat map creation and retrieval
- Route stop ordering
- Filter combinations

### Performance Tests
- Search with 10,000+ trips
- Concurrent seat availability requests
- Cache effectiveness monitoring

---

## Success Criteria

✅ Admins can create routes with multiple stops
✅ Admins can configure custom seat layouts
✅ Users can search with 7+ filter options
✅ Search results include pricing tiers
✅ Seat types have automatic price multipliers
✅ Database queries use proper indexes
✅ Popular searches are cached
✅ All constants are type-safe

---

## Timeline

**Total Estimated Time:** ~15 hours

**Day 1 (4 hours):**
- Constants package
- Route stop model & migration
- Seat model enhancement

**Day 2 (5 hours):**
- Repositories implementation
- Database indexes
- Cache service

**Day 3 (4 hours):**
- Service layer
- Advanced search logic
- Handlers

**Day 4 (2 hours):**
- Testing
- Documentation
- Code review

---

## Dependencies

### Technical
- Redis for caching
- PostgreSQL with PostGIS (for GPS coordinates)
- Existing microservices architecture

### Business Logic
- Seat pricing rules
- Route stop validation rules
- Search ranking algorithm

---

## Risks & Mitigation

**Risk:** Complex search queries may be slow
- **Mitigation:** Database indexes, caching, pagination

**Risk:** Seat map configuration may be complex for admins
- **Mitigation:** Bulk creation, templates, visual preview

**Risk:** Cache invalidation complexity
- **Mitigation:** Clear invalidation rules, short TTLs

---

## Future Enhancements (Week 3+)

- Real-time seat availability
- Dynamic pricing based on demand
- Route optimization algorithms
- Mobile app integration
- Payment gateway integration

---

## Notes

This plan builds upon Week 1's authentication and authorization foundation. All features are designed to be scalable and maintainable, following the established microservices architecture and clean code principles.
