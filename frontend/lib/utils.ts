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
 * Format a Date object to HTML5 date input format (yyyy-MM-dd)
 * @param date - The date to format
 * @returns Date string in yyyy-MM-dd format
 * @example formatDateForInput(new Date('2024-12-25')) // returns "2024-12-25"
 */
export function formatDateForInput(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
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

import { DisplayValue, type ConstantDisplay } from "@/lib/types/trip";

/**
 * Type guard to check if a value is a DisplayValue object
 * @param value - Value to check
 * @returns True if value is a DisplayValue object
 */
export function isDisplayValue(value: unknown): value is DisplayValue<string> {
  return (
    typeof value === "object" &&
    value !== null &&
    "value" in value &&
    "display_name" in value
  );
}

/**
 * Extract display name from DisplayValue object or return string as-is
 * Provides backward compatibility with raw string values
 * @param value - DisplayValue object or raw string
 * @param fallback - Fallback string if extraction fails
 * @returns Extracted display name or fallback
 */
export function getDisplayName(
  value: string | ConstantDisplay | undefined | null,
  fallback: string = "N/A",
): string {
  if (!value) return fallback;
  if (typeof value === "string") return value;
  if (isDisplayValue(value)) return value.display_name;
  return fallback;
}

/**
 * Extract raw value from DisplayValue object or return string as-is
 * @param value - DisplayValue object or raw string
 * @returns Extracted value or original string
 */
export function getValue(
  value: string | ConstantDisplay | undefined | null,
): string {
  if (!value) return "";
  if (typeof value === "string") return value;
  if (isDisplayValue(value)) return value.value;
  return "";
}

import { ApiTripItem, TripDetail } from "@/lib/types/trip";
import { getAmenityDisplay, getTripStatusDisplay } from "@/lib/constants/trip";

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

  // Map raw amenity strings to display names using constants
  const bus_amenities = apiTrip.bus.amenities.map((amenity) =>
    getAmenityDisplay(amenity),
  );

  return {
    id: apiTrip.id,
    route_id: apiTrip.route_id,
    bus_id: apiTrip.bus_id,
    departure_time: apiTrip.departure_time,
    arrival_time: apiTrip.arrival_time,
    base_price: apiTrip.base_price,
    status: getTripStatusDisplay(apiTrip.status), // Map raw status to display name
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
