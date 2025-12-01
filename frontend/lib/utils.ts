import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Format a Date object to Vietnamese date format (dd/MM/yyyy) for API calls
 * @param date - The date to format
 * @returns Date string in dd/MM/yyyy format
 * @example formatDateForApi(new Date('2024-12-25')) // returns "25/12/2024"
 */
export function formatDateForApi(date: Date): string {
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}/${month}/${year}`;
}

/**
 * Parse a Vietnamese date format (dd/MM/yyyy) string to a Date object
 * @param dateString - The date string in dd/MM/yyyy format
 * @returns Date object or null if invalid
 * @example parseDateFromVnFormat("25/12/2024") // returns Date object for Dec 25, 2024
 */
export function parseDateFromVnFormat(dateString: string): Date | null {
  const parts = dateString.split("/");
  if (parts.length !== 3) return null;

  const day = parseInt(parts[0], 10);
  const month = parseInt(parts[1], 10);
  const year = parseInt(parts[2], 10);

  if (isNaN(day) || isNaN(month) || isNaN(year)) return null;

  const date = new Date(year, month - 1, day);

  // Validate the date
  if (
    date.getDate() !== day ||
    date.getMonth() !== month - 1 ||
    date.getFullYear() !== year
  ) {
    return null;
  }

  return date;
}

import { ApiTripItem, TripDetail } from "@/lib/types/trip";

/**
 * Transform API trip item to internal TripDetail format
 * Maintains backward compatibility with existing components
 */
export function transformApiTripToTripDetail(apiTrip: ApiTripItem): TripDetail {
  // Calculate duration string (e.g., "30h 0m")
  const durationMinutes = apiTrip.route.duration_minutes;
  const hours = Math.floor(durationMinutes / 60);
  const minutes = durationMinutes % 60;
  const duration = `${hours}h ${minutes}m`;

  // Extract amenity display names
  const bus_amenities = apiTrip.bus.amenities.map(
    (amenity) => amenity.display_name,
  );

  return {
    id: apiTrip.id,
    route_id: apiTrip.route_id,
    bus_id: apiTrip.bus_id,
    departure_time: apiTrip.departure_time,
    arrival_time: apiTrip.arrival_time,
    base_price: apiTrip.base_price,
    status: apiTrip.status.display_name,
    available_seats: apiTrip.available_seats,
    total_seats: apiTrip.total_seats,
    duration,
    origin: apiTrip.route.origin,
    destination: apiTrip.route.destination,
    distance_km: apiTrip.route.distance_km,
    bus_model: apiTrip.bus.model,
    bus_plate_number: "", // Not provided in new API
    bus_amenities,
    operator_id: "", // Not provided in new API, would need separate operator lookup
    operator_name: "", // Not provided in new API
  };
}
