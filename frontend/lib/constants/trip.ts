/**
 * Trip-related constants for frontend display
 * These map the raw string values from backend to user-friendly display names
 */

// Seat Type Constants
export const SEAT_TYPES = {
  standard: {
    value: "standard",
    displayName: "Ghế thường",
    priceMultiplier: 1.0,
  },
  vip: {
    value: "vip",
    displayName: "Ghế VIP",
    priceMultiplier: 1.2,
  },
  sleeper: {
    value: "sleeper",
    displayName: "Giường nằm",
    priceMultiplier: 1.5,
  },
} as const;

export type SeatType = keyof typeof SEAT_TYPES;

export function getSeatTypeDisplay(seatType: string): string {
  const type = SEAT_TYPES[seatType as SeatType];
  return type?.displayName || seatType;
}

export function getSeatTypePriceMultiplier(seatType: string): number {
  const type = SEAT_TYPES[seatType as SeatType];
  return type?.priceMultiplier || 1.0;
}

// Amenity Constants
export const AMENITIES = {
  wifi: {
    value: "wifi",
    displayName: "WiFi",
    icon: "📶",
  },
  ac: {
    value: "ac",
    displayName: "Điều hòa",
    icon: "❄️",
  },
  toilet: {
    value: "toilet",
    displayName: "Nhà vệ sinh",
    icon: "🚻",
  },
  tv: {
    value: "tv",
    displayName: "TV",
    icon: "📺",
  },
  water: {
    value: "water",
    displayName: "Nước uống",
    icon: "💧",
  },
  blanket: {
    value: "blanket",
    displayName: "Chăn",
    icon: "🛏️",
  },
  usb_charger: {
    value: "usb_charger",
    displayName: "Sạc USB",
    icon: "🔌",
  },
  snack: {
    value: "snack",
    displayName: "Đồ ăn nhẹ",
    icon: "🍪",
  },
} as const;

export type Amenity = keyof typeof AMENITIES;

export function getAmenityDisplay(amenity: string): string {
  const item = AMENITIES[amenity as Amenity];
  return item?.displayName || amenity;
}

export function getAmenityIcon(amenity: string): string {
  const item = AMENITIES[amenity as Amenity];
  return item?.icon || "";
}

// Bus Type Constants
export const BUS_TYPES = {
  standard: {
    value: "standard",
    displayName: "Xe thường",
  },
  limousine: {
    value: "limousine",
    displayName: "Limousine",
  },
  sleeper: {
    value: "sleeper",
    displayName: "Giường nằm",
  },
} as const;

export type BusType = keyof typeof BUS_TYPES;

export function getBusTypeDisplay(busType: string): string {
  const type = BUS_TYPES[busType as BusType];
  return type?.displayName || busType;
}

// Trip Status Constants
export const TRIP_STATUSES = {
  scheduled: {
    value: "scheduled",
    displayName: "Đã lên lịch",
    color: "blue",
    variant: "default" as const,
  },
  in_progress: {
    value: "in_progress",
    displayName: "Đang chạy",
    color: "green",
    variant: "default" as const,
  },
  completed: {
    value: "completed",
    displayName: "Hoàn thành",
    color: "gray",
    variant: "secondary" as const,
  },
  cancelled: {
    value: "cancelled",
    displayName: "Đã hủy",
    color: "red",
    variant: "destructive" as const,
  },
} as const;

export type TripStatus = keyof typeof TRIP_STATUSES;

export function getTripStatusDisplay(status: string): string {
  const tripStatus = TRIP_STATUSES[status as TripStatus];
  return tripStatus?.displayName || status;
}

export function getTripStatusVariant(
  status: string,
): "default" | "secondary" | "destructive" {
  const tripStatus = TRIP_STATUSES[status as TripStatus];
  return tripStatus?.variant || "default";
}

// Stop Type Constants
export const STOP_TYPES = {
  pickup: {
    value: "pickup",
    displayName: "Điểm đón",
  },
  dropoff: {
    value: "dropoff",
    displayName: "Điểm trả",
  },
  rest: {
    value: "rest",
    displayName: "Điểm nghỉ",
  },
} as const;

export type StopType = keyof typeof STOP_TYPES;

export function getStopTypeDisplay(stopType: string): string {
  const type = STOP_TYPES[stopType as StopType];
  return type?.displayName || stopType;
}

// Helper to get all amenity options for filters
export function getAllAmenityOptions() {
  return Object.values(AMENITIES);
}

// Helper to get all seat type options for filters
export function getAllSeatTypeOptions() {
  return Object.values(SEAT_TYPES);
}

// Helper to get all bus type options
export function getAllBusTypeOptions() {
  return Object.values(BUS_TYPES);
}
