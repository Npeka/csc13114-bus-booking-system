"use client";

import { useState, Suspense, useMemo } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { type Trip } from "@/components/trips/trip-card";
import { TripFilters, type Filters } from "@/components/trips/trip-filters";
import { TripSearchForm } from "@/components/search/trip-search-form";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { PaginationWithLinks } from "@/components/ui/pagination-with-links";
import { Filter } from "lucide-react";
import { searchTrips } from "@/lib/api/trip/trip-service";
import type { ApiTripItem, TripSearchParams } from "@/lib/types/trip";
import { parseDateFromVnFormat } from "@/lib/utils";
import { TripSummaryHeader } from "./_components/trip-summary-header";
import { TripResults } from "./_components/trip-results";

function TripsContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [sortBy, setSortBy] = useState<"price" | "departure" | "duration">(
    "price",
  );
  const pageSize = 20;

  const origin = searchParams.get("from") || "";
  const destination = searchParams.get("to") || "";
  const date = searchParams.get("date") || "";
  const passengers = parseInt(searchParams.get("passengers") || "1", 10);

  const currentPage = Number(searchParams.get("page")) || 1;

  const createPageURL = (pageNumber: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("page", pageNumber.toString());
    return `/trips?${params.toString()}`;
  };

  const [filters, setFilters] = useState<Filters>({
    priceRange: [0, 1000000],
    departureTime: [],
    busTypes: [],
    amenities: [],
    operators: [],
  });

  // Map departure time slots to time ranges
  const getTimeRange = (slots: string[]) => {
    if (slots.length === 0) return { min: undefined, max: undefined };

    const timeRanges: Record<string, { min: string; max: string }> = {
      morning: { min: "00:00", max: "06:00" },
      daytime: { min: "06:00", max: "12:00" },
      afternoon: { min: "12:00", max: "18:00" },
      evening: { min: "18:00", max: "24:00" },
    };

    const mins = slots.map((s) => timeRanges[s]?.min).filter(Boolean);
    const maxs = slots.map((s) => timeRanges[s]?.max).filter(Boolean);

    return {
      min: mins.length > 0 ? mins.sort()[0] : undefined,
      max: maxs.length > 0 ? maxs.sort().reverse()[0] : undefined,
    };
  };

  // Build search params for API
  const searchParams_api: TripSearchParams = useMemo(() => {
    // Convert date from dd/MM/yyyy to ISO date range
    let departureStart: string | undefined;
    let departureEnd: string | undefined;

    if (date) {
      const parsedDate = parseDateFromVnFormat(date);
      if (parsedDate) {
        // Start of day
        const startOfDay = new Date(parsedDate);
        startOfDay.setHours(0, 0, 0, 0);
        departureStart = startOfDay.toISOString();

        // End of day
        const endOfDay = new Date(parsedDate);
        endOfDay.setHours(23, 59, 59, 999);
        departureEnd = endOfDay.toISOString();
      }
    }

    const params: TripSearchParams = {
      origin: origin || undefined,
      destination: destination || undefined,
      departure_time_start: departureStart,
      departure_time_end: departureEnd,
      passengers, // Client-side only
      page: currentPage,
      page_size: pageSize,
    };

    // Add sorting
    if (sortBy === "price") {
      params.sort_by = "price";
      params.sort_order = "asc";
    } else if (sortBy === "departure") {
      params.sort_by = "departure_time";
      params.sort_order = "asc";
    }

    // Add price filter
    if (filters.priceRange[0] > 0) {
      params.min_price = filters.priceRange[0];
    }
    if (filters.priceRange[1] < 1000000) {
      params.max_price = filters.priceRange[1];
    }

    // Add time range filter
    if (filters.departureTime.length > 0) {
      const timeRange = getTimeRange(filters.departureTime);
      if (timeRange.min) params.departure_time_min = timeRange.min;
      if (timeRange.max) params.departure_time_max = timeRange.max;
    }

    // Add amenities filter
    if (filters.amenities.length > 0) {
      params.amenities = filters.amenities;
    }

    // Add bus type filter (map Vietnamese names to search terms)
    if (filters.busTypes.length > 0) {
      // Use first bus type for now (can be enhanced to support multiple)
      const busTypeMap: Record<string, string> = {
        "Ghế ngồi": "seat",
        "Giường nằm": "bed",
        Limousine: "limousine",
        "Cabin đôi": "cabin",
      };
      const mappedType = busTypeMap[filters.busTypes[0]] || filters.busTypes[0];
      params.bus_type = mappedType;
    }

    return params;
  }, [
    origin,
    destination,
    date,
    passengers,
    currentPage,
    sortBy,
    filters.priceRange,
    filters.departureTime,
    filters.amenities,
    filters.busTypes,
  ]);

  // Fetch trips from API
  const {
    data: searchResponse,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["trips", searchParams_api],
    queryFn: () => searchTrips(searchParams_api),
    enabled: !!origin && !!destination && !!date,
  });

  // Convert API TripDetail to Trip format for TripCard
  const trips: Trip[] = useMemo(() => {
    if (!searchResponse?.trips) return [];

    return searchResponse.trips.map((apiTrip: ApiTripItem) => {
      const departureDate = new Date(apiTrip.departure_time);
      const arrivalDate = new Date(apiTrip.arrival_time);
      const durationMs = arrivalDate.getTime() - departureDate.getTime();
      const hours = Math.floor(durationMs / (1000 * 60 * 60));
      const minutes = Math.floor((durationMs % (1000 * 60 * 60)) / (1000 * 60));

      // Map bus type from backend value
      let busType = "Ghế ngồi";
      if (apiTrip.bus.bus_type === "sleeper") {
        busType = "Giường nằm";
      } else if (apiTrip.bus.bus_type === "limousine") {
        busType = "Limousine";
      }

      // Map amenities from raw strings to display names
      const amenityNames = (apiTrip.bus.amenities || []).map((amenity) => {
        const mapping: Record<string, string> = {
          wifi: "WiFi",
          ac: "Điều hòa",
          toilet: "Toilet",
          tv: "TV",
          water: "Nước uống",
          blanket: "Chăn",
          usb_charger: "Sạc USB",
          snack: "Snack",
        };
        return mapping[amenity] || amenity;
      });

      return {
        id: apiTrip.id,
        operator: "Nhà xe", // operator info not in search API
        operatorRating: 4.5, // Default rating, can be enhanced later
        departureTime: format(departureDate, "HH:mm"),
        arrivalTime: format(arrivalDate, "HH:mm"),
        duration: `${hours}h ${minutes}m`,
        origin: apiTrip.route.origin,
        destination: apiTrip.route.destination,
        price: apiTrip.base_price,
        availableSeats: apiTrip.available_seats,
        busType,
        amenities: amenityNames,
      };
    });
  }, [searchResponse]);

  const handleClearFilters = () => {
    setFilters({
      priceRange: [0, 1000000],
      departureTime: [],
      busTypes: [],
      amenities: [],
      operators: [],
    });
  };

  const handleSelectTrip = (tripId: string) => {
    router.push(`/trips/${tripId}`);
  };

  // Apply client-side filters (for bus types and amenities that aren't in API yet)
  const filteredTrips = trips.filter((trip) => {
    if (
      filters.busTypes.length > 0 &&
      !filters.busTypes.includes(trip.busType)
    ) {
      return false;
    }
    if (filters.amenities.length > 0) {
      const hasAllAmenities = filters.amenities.every((amenity) =>
        trip.amenities.includes(amenity),
      );
      if (!hasAllAmenities) return false;
    }
    // Departure time filtering will be handled by backend in future
    return true;
  });

  // Client-side sorting for duration (API handles price and departure)
  const sortedTrips =
    sortBy === "duration"
      ? [...filteredTrips].sort((a, b) => {
          const aDuration = parseInt(a.duration.replace(/[^\d]/g, "")) || 0;
          const bDuration = parseInt(b.duration.replace(/[^\d]/g, "")) || 0;
          return aDuration - bDuration;
        })
      : filteredTrips;

  const activeFiltersCount =
    (filters.priceRange[0] !== 0 || filters.priceRange[1] !== 1000000 ? 1 : 0) +
    filters.departureTime.length +
    filters.busTypes.length +
    filters.amenities.length;

  const handleToggleSort = () => {
    setSortBy((prev) =>
      prev === "price"
        ? "departure"
        : prev === "departure"
          ? "duration"
          : "price",
    );
  };

  return (
    <div className="min-h-screen">
      {/* Update Search Form */}
      <div className="border-b py-6">
        <div className="container">
          <TripSearchForm />
        </div>
      </div>

      {/* Main Content */}
      <div className="container py-8">
        <div className="grid gap-6 lg:grid-cols-[280px_1fr]">
          {/* Desktop Filters */}
          <aside className="hidden lg:block">
            <div className="sticky top-20">
              <TripFilters
                filters={filters}
                onFiltersChange={setFilters}
                onClearFilters={handleClearFilters}
              />
            </div>
          </aside>

          {/* Mobile Filters */}
          <div className="lg:hidden">
            <Sheet>
              <SheetTrigger asChild>
                <Button variant="outline" className="w-full">
                  <Filter className="mr-2 h-4 w-4" />
                  Bộ lọc
                  {activeFiltersCount > 0 && (
                    <Badge variant="secondary" className="ml-2">
                      {activeFiltersCount}
                    </Badge>
                  )}
                </Button>
              </SheetTrigger>
              <SheetContent side="left" className="w-80 overflow-y-auto">
                <TripFilters
                  filters={filters}
                  onFiltersChange={setFilters}
                  onClearFilters={handleClearFilters}
                />
              </SheetContent>
            </Sheet>
          </div>

          {/* Trip List */}
          <div className="space-y-4">
            <TripSummaryHeader
              origin={origin}
              destination={destination}
              date={date}
              passengers={passengers.toString()}
              resultsCount={searchResponse?.total || sortedTrips.length}
              sortBy={sortBy}
              onToggleSort={handleToggleSort}
            />
            <TripResults
              loading={isLoading}
              trips={sortedTrips}
              onSelect={handleSelectTrip}
              onClearFilters={handleClearFilters}
              error={error}
            />
            {searchResponse && searchResponse.total_pages > 1 && (
              <div className="mt-6">
                <PaginationWithLinks
                  page={currentPage}
                  totalPages={searchResponse.total_pages}
                  createPageURL={createPageURL}
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default function TripsPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <TripsContent />
    </Suspense>
  );
}
