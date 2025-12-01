// Trip search request parameters
export interface TripSearchParams {
  origin: string;
  destination: string;
  date: string; // Vietnamese date format: dd/MM/yyyy
  passengers: number;
  seat_type?: "standard" | "vip" | "sleeper";
  price_min?: number;
  price_max?: number;
  departure_time_min?: string; // Format: "HH:MM" (e.g., "06:00")
  departure_time_max?: string; // Format: "HH:MM" (e.g., "18:00")
  amenities?: string[]; // Filter by bus amenities
  bus_type?: string; // Filter by bus type/model
  operator_id?: string;
  sort_by?: "price" | "departure_time" | "arrival_time";
  sort_order?: "asc" | "desc";
  page?: number;
  limit?: number;
}

/**
 * Generic type for API enum values with display names
 */
export interface DisplayValue<T = string> {
  value: T;
  display_name: string;
}

/**
 * Route information from API response
 */
export interface ApiTripRoute {
  id: string;
  origin: string;
  destination: string;
  distance_km: number;
  duration_minutes: number;
}

/**
 * Amenity from API response
 */
export interface ApiAmenity {
  value: string;
  display_name: string;
}

/**
 * Bus information from API response
 */
export interface ApiBusInfo {
  id: string;
  model: string;
  bus_type: DisplayValue<string>;
  total_seats: number;
  amenities: ApiAmenity[];
}

/**
 * Trip item from API search response
 */
export interface ApiTripItem {
  id: string;
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime with timezone
  arrival_time: string; // ISO datetime with timezone
  base_price: number;
  status: DisplayValue<string>;
  available_seats: number;
  total_seats: number;
  route: ApiTripRoute;
  bus: ApiBusInfo;
}

/**
 * Pagination metadata from API
 */
export interface ApiPaginationMeta {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

/**
 * Generic paginated response structure from backend
 */
export interface ApiPaginatedResponse<T> {
  data: T[];
  meta: ApiPaginationMeta;
}

/**
 * Complete API search response structure
 */
export interface ApiTripSearchResponse {
  data: ApiTripItem[];
  meta: ApiPaginationMeta;
}

// Trip detail from search response
export interface TripDetail {
  id: string;
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime string
  arrival_time: string; // ISO datetime string
  base_price: number;
  status: string;
  available_seats: number;
  total_seats: number;
  duration: string;
  origin: string;
  destination: string;
  distance_km: number;
  bus_model: string;
  bus_plate_number: string;
  bus_amenities: string[];
  operator_id: string;
  operator_name: string;
}

// Trip search response
export interface TripSearchResponse {
  trips: TripDetail[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// Route types
export interface Route {
  id: string;
  origin: string;
  destination: string;
  distance_km: number;
  estimated_minutes: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  route_stops?: RouteStop[]; // Optional as it might not be returned in list views
}

// Bus Seat type from API
export interface BusSeat {
  id: string;
  seat_number: string;
  row: number;
  column: number;
  floor: number;
  seat_type: "standard" | "vip" | "sleeper";
  price_multiplier: number;
  is_available: boolean;
}

// Seat layout editor types
export type SeatCellType = "seat" | "empty" | "blocked" | "driver";

export interface SeatLayoutCell {
  id?: string; // Seat ID if it's a seat
  type: SeatCellType;
  seatType?: "standard" | "vip" | "sleeper";
  seatNumber?: string;
  priceMultiplier?: number;
  isAvailable?: boolean;
  row: number;
  column: number;
  floor: number;
}

export interface SeatLayoutFloor {
  floor: number;
  rows: number;
  cols: number;
  cells: SeatLayoutCell[][];
}

export interface SeatLayoutConfig {
  busId: string;
  floors: SeatLayoutFloor[];
}

export interface CreateSeatRequest {
  bus_id: string;
  seat_number: string;
  row: number;
  column: number;
  seat_type: "standard" | "vip" | "sleeper";
  price_multiplier?: number;
  floor: number;
}

export interface BulkCreateSeatsRequest {
  bus_id: string;
  seats: CreateSeatRequest[];
}

// Bus types
export interface Bus {
  id: string;
  plate_number: string;
  model: string;
  seat_capacity: number;
  amenities: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
  seats?: BusSeat[]; // Optional as it might not be returned in list views
}

// Trip types
export interface Trip {
  id: string;
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime string
  arrival_time: string; // ISO datetime string
  base_price: number;
  status: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  route?: Route; // Populated when fetching with preload
  bus?: Bus; // Populated when fetching with preload
}

// Seat types
export interface Seat {
  id: string;
  bus_id: string;
  seat_code: string;
  seat_type: "standard" | "vip" | "sleeper";
  is_active: boolean;
}

// Seat detail from trip seats response
export interface SeatDetail {
  id: string;
  seat_code: string;
  seat_type: string;
  is_booked: boolean;
  is_locked: boolean;
  price: number;
}

export interface SeatAvailabilityResponse {
  trip_id: string;
  available_seats: number;
  total_seats: number;
  seats: SeatDetail[];
}

// Route Stop types - matches backend model
export interface RouteStop {
  id: string;
  route_id: string;
  stop_order: number;
  stop_type: string; // "pickup" | "dropoff" | "both"
  location: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  offset_minutes: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateRouteStopRequest {
  stop_order: number;
  stop_type: string; // "pickup" | "dropoff" | "both"
  location: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  offset_minutes: number;
}

export interface UpdateRouteStopRequest {
  stop_order?: number;
  stop_type?: string;
  location?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  offset_minutes?: number;
  is_active?: boolean;
}
