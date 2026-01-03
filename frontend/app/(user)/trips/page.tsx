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
import { TripPagination } from "./_components/trip-pagination";

function TripsContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  // Read sortBy and sortOrder from URL params
  const sortBy =
    (searchParams.get("sort") as "price" | "departure" | "duration") || "price";
  const sortOrder = (searchParams.get("order") as "asc" | "desc") || "asc";
  const pageSize = 20;

  const origin = (searchParams.get("from") || "").replace(/\+/g, " ");
  const destination = (searchParams.get("to") || "").replace(/\+/g, " ");
  const date = (searchParams.get("date") || "").replace(/\+/g, " ");
  const passengers = parseInt(searchParams.get("passengers") || "1", 10);

  const currentPage = Number(searchParams.get("page")) || 1;

  const createPageURL = (pageNumber: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("page", pageNumber.toString());
    return `/trips?${params.toString()}`;
  };

  // Read filters from URL params
  const priceMin = parseInt(searchParams.get("priceMin") || "0", 10);
  const priceMax = parseInt(searchParams.get("priceMax") || "1000000", 10);
  const departureTimeSlots =
    searchParams.get("departureTime")?.split(",").filter(Boolean) || [];
  const busTypes =
    searchParams.get("busTypes")?.split(",").filter(Boolean) || [];
  const amenitiesFromUrl =
    searchParams.get("amenities")?.split(",").filter(Boolean) || [];

  const filters: Filters = {
    priceRange: [priceMin, priceMax],
    departureTime: departureTimeSlots,
    busTypes: busTypes,
    amenities: amenitiesFromUrl,
    operators: [],
  };

  // Map departure time slots (start times) to time ranges
  const getTimeRange = (slots: string[]) => {
    if (slots.length === 0) return { min: undefined, max: undefined };

    // Map start_time to end_time
    // Each slot is 6 hours long
    const endTimes: Record<string, string> = {
      "00:00": "06:00",
      "06:00": "12:00",
      "12:00": "18:00",
      "18:00": "24:00",
    };

    const validSlots = slots.filter((s) => endTimes[s]);
    if (validSlots.length === 0) return { min: undefined, max: undefined };

    const mins = validSlots.sort();
    const maxs = validSlots.map((s) => endTimes[s]).sort();

    return {
      min: mins[0],
      max: maxs[maxs.length - 1],
    };
  };

  // Build search params for API
  const searchParams_api: TripSearchParams = useMemo(() => {
    // Convert date from dd/MM/yyyy to ISO date range
    let departureStart: string | undefined;
    let departureEnd: string | undefined;
    let parsedDate: Date | null = null;

    if (date) {
      parsedDate = parseDateFromVnFormat(date);
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
    } else if (sortBy === "departure") {
      params.sort_by = "departure_time";
    } else if (sortBy === "duration") {
      params.sort_by = "duration";
    }
    params.sort_order = sortOrder;

    // Add price filter
    if (filters.priceRange[0] > 0) {
      params.min_price = filters.priceRange[0];
    }
    if (filters.priceRange[1] < 1000000) {
      params.max_price = filters.priceRange[1];
    }

    // Add time range filter (merge date with time slots)
    // If time slots are selected, constrain the start/end times
    // Otherwise, use the full day range calculated above
    if (filters.departureTime.length > 0 && parsedDate) {
      const timeRange = getTimeRange(filters.departureTime);

      if (timeRange.min) {
        const [h, m] = timeRange.min.split(":").map(Number);
        const start = new Date(parsedDate);
        start.setHours(h, m, 0, 0);
        params.departure_time_start = start.toISOString();
      }

      if (timeRange.max) {
        const [h, m] = timeRange.max.split(":").map(Number);
        const end = new Date(parsedDate);
        end.setHours(h, m, 59, 999);
        params.departure_time_end = end.toISOString();
      }
    }

    // Add amenities filter - filter values are already raw codes (wifi, ac, etc)
    if (filters.amenities.length > 0) {
      params.amenities = filters.amenities;
    }

    // Add bus/seat type filter
    // Frontend filter combines Bus Type and Seat Type concepts
    // We map these to the backend's seat_types[] filter
    if (filters.busTypes.length > 0) {
      const seatTypeMap: Record<string, string> = {
        "Ghế thường": "standard", // Changed from "Ghế ngồi" to match UI likely
        "Ghế ngồi": "standard",
        "Giường nằm": "sleeper",
        "Ghế VIP": "vip",
        Limousine: "vip", // Map Limousine to VIP seat for now
      };

      // Filter out any unmapped types to avoid sending invalid values
      const seatTypes = filters.busTypes
        .map((t) => seatTypeMap[t] || t.toLowerCase())
        .filter((t) => ["standard", "vip", "sleeper"].includes(t));

      if (seatTypes.length > 0) {
        params.seat_types = seatTypes;
      }
    }

    return params;
  }, [
    origin,
    destination,
    date,
    passengers,
    currentPage,
    sortBy,
    sortOrder,
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
    refetchInterval: 2000, // Poll every 2 seconds for realtime seat updates
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
          charging: "Sạc USB",
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

  const handleFiltersChange = (newFilters: Filters) => {
    const params = new URLSearchParams(searchParams.toString());

    // Update price range
    if (newFilters.priceRange[0] > 0) {
      params.set("priceMin", newFilters.priceRange[0].toString());
    } else {
      params.delete("priceMin");
    }
    if (newFilters.priceRange[1] < 1000000) {
      params.set("priceMax", newFilters.priceRange[1].toString());
    } else {
      params.delete("priceMax");
    }

    // Update departure time
    if (newFilters.departureTime.length > 0) {
      params.set("departureTime", newFilters.departureTime.join(","));
    } else {
      params.delete("departureTime");
    }

    // Update bus types
    if (newFilters.busTypes.length > 0) {
      params.set("busTypes", newFilters.busTypes.join(","));
    } else {
      params.delete("busTypes");
    }

    // Update amenities
    if (newFilters.amenities.length > 0) {
      params.set("amenities", newFilters.amenities.join(","));
    } else {
      params.delete("amenities");
    }

    // Reset to page 1 when filters change
    params.set("page", "1");

    router.push(`/trips?${params.toString()}`, { scroll: false });
  };

  const handleClearFilters = () => {
    const params = new URLSearchParams(searchParams.toString());
    params.delete("priceMin");
    params.delete("priceMax");
    params.delete("departureTime");
    params.delete("busTypes");
    params.delete("amenities");
    params.set("page", "1");
    router.push(`/trips?${params.toString()}`, { scroll: false });
  };

  const handleSelectTrip = (tripId: string) => {
    router.push(`/trips/${tripId}`);
  };

  // Remove client-side filtering - backend handles all filters
  const filteredTrips = trips;

  // Client-side sorting removed - backend handles all sorting including duration
  const sortedTrips = filteredTrips;

  const activeFiltersCount =
    (filters.priceRange[0] !== 0 || filters.priceRange[1] !== 1000000 ? 1 : 0) +
    filters.departureTime.length +
    filters.busTypes.length +
    filters.amenities.length;

  const handleSortChange = (field: string, order: "asc" | "desc") => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("sort", field);
    params.set("order", order);
    router.push(`/trips?${params.toString()}`, { scroll: false });
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
                onFiltersChange={handleFiltersChange}
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
                  onFiltersChange={handleFiltersChange}
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
              sortOrder={sortOrder}
              onSortChange={handleSortChange}
            />
            <TripResults
              loading={isLoading}
              trips={sortedTrips}
              onSelect={handleSelectTrip}
              onClearFilters={handleClearFilters}
              error={error}
            />
            {searchResponse && (
              <TripPagination
                currentPage={currentPage}
                totalPages={searchResponse.total_pages}
                onPageCreateURL={createPageURL}
              />
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
