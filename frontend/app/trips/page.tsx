"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { TripCard, TripCardSkeleton, type Trip } from "@/components/trips/trip-card";
import { TripFilters, type Filters } from "@/components/trips/trip-filters";
import { TripSearchForm } from "@/components/search/trip-search-form";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Filter, ArrowUpDown } from "lucide-react";

function TripsContent() {
  const searchParams = useSearchParams();
  const [trips, setTrips] = useState<Trip[]>([]);
  const [loading, setLoading] = useState(true);
  const [sortBy, setSortBy] = useState<"price" | "departure" | "duration">("price");

  const origin = searchParams.get("from") || "";
  const destination = searchParams.get("to") || "";
  const date = searchParams.get("date") || "";
  const passengers = searchParams.get("passengers") || "1";

  const [filters, setFilters] = useState<Filters>({
    priceRange: [0, 1000000],
    departureTime: [],
    busTypes: [],
    amenities: [],
    operators: [],
  });

  useEffect(() => {
    // Simulate API call
    setLoading(true);
    setTimeout(() => {
      setTrips(mockTrips);
      setLoading(false);
    }, 1000);
  }, [origin, destination, date]);

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
    // Navigate to seat selection
    window.location.href = `/trips/${tripId}/seats`;
  };

  const filteredTrips = trips.filter((trip) => {
    // Apply filters
    if (trip.price < filters.priceRange[0] || trip.price > filters.priceRange[1]) {
      return false;
    }
    if (filters.busTypes.length > 0 && !filters.busTypes.includes(trip.busType)) {
      return false;
    }
    if (filters.amenities.length > 0) {
      const hasAllAmenities = filters.amenities.every((amenity) =>
        trip.amenities.includes(amenity)
      );
      if (!hasAllAmenities) return false;
    }
    // Add departure time filtering logic
    return true;
  });

  const sortedTrips = [...filteredTrips].sort((a, b) => {
    switch (sortBy) {
      case "price":
        return a.price - b.price;
      case "departure":
        return a.departureTime.localeCompare(b.departureTime);
      case "duration":
        return a.duration.localeCompare(b.duration);
      default:
        return 0;
    }
  });

  const activeFiltersCount =
    (filters.priceRange[0] !== 0 || filters.priceRange[1] !== 1000000 ? 1 : 0) +
    filters.departureTime.length +
    filters.busTypes.length +
    filters.amenities.length;

  return (
    <div className="min-h-screen bg-neutral-50">
      {/* Search Summary Bar */}
      <div className="bg-white border-b">
        <div className="container py-4">
          <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
            <div>
              <h1 className="text-2xl font-bold">
                {origin} → {destination}
              </h1>
              <p className="text-sm text-muted-foreground">
                {date} • {passengers} hành khách • {sortedTrips.length} chuyến xe
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setSortBy(sortBy === "price" ? "departure" : "price")}
              >
                <ArrowUpDown className="mr-2 h-4 w-4" />
                Sắp xếp: {sortBy === "price" ? "Giá" : sortBy === "departure" ? "Giờ đi" : "Thời gian"}
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Update Search Form */}
      <div className="bg-white border-b py-6">
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
            {loading ? (
              <>
                {[...Array(5)].map((_, i) => (
                  <TripCardSkeleton key={i} />
                ))}
              </>
            ) : sortedTrips.length > 0 ? (
              sortedTrips.map((trip) => (
                <TripCard key={trip.id} trip={trip} onSelect={handleSelectTrip} />
              ))
            ) : (
              <div className="text-center py-12">
                <p className="text-lg text-muted-foreground">
                  Không tìm thấy chuyến xe phù hợp
                </p>
                <Button
                  variant="outline"
                  className="mt-4"
                  onClick={handleClearFilters}
                >
                  Xóa bộ lọc
                </Button>
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

// Mock data
const mockTrips: Trip[] = [
  {
    id: "1",
    operator: "Phương Trang FUTA Bus Lines",
    operatorRating: 4.8,
    departureTime: "06:00",
    arrivalTime: "14:30",
    duration: "8h 30m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    price: 180000,
    availableSeats: 15,
    busType: "Giường nằm",
    amenities: ["WiFi", "Điều hòa", "Nước uống", "Sạc điện thoại"],
  },
  {
    id: "2",
    operator: "Mai Linh Express",
    operatorRating: 4.6,
    departureTime: "07:30",
    arrivalTime: "16:00",
    duration: "8h 30m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    price: 165000,
    availableSeats: 8,
    busType: "Ghế ngồi",
    amenities: ["Điều hòa", "Nước uống"],
  },
  {
    id: "3",
    operator: "Thành Bưởi Limousine",
    operatorRating: 4.9,
    departureTime: "08:00",
    arrivalTime: "16:15",
    duration: "8h 15m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    price: 250000,
    availableSeats: 12,
    busType: "Limousine",
    amenities: ["WiFi", "Điều hòa", "Nước uống", "Sạc điện thoại", "TV"],
  },
  {
    id: "4",
    operator: "Kumho Samco",
    operatorRating: 4.5,
    departureTime: "09:30",
    arrivalTime: "18:00",
    duration: "8h 30m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    price: 175000,
    availableSeats: 20,
    busType: "Giường nằm",
    amenities: ["WiFi", "Điều hòa", "Nước uống"],
  },
  {
    id: "5",
    operator: "Hanh Cafe",
    operatorRating: 4.7,
    departureTime: "22:00",
    arrivalTime: "06:30",
    duration: "8h 30m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    price: 190000,
    availableSeats: 10,
    busType: "Giường nằm",
    amenities: ["WiFi", "Điều hòa", "Nước uống", "Sạc điện thoại", "Toilet"],
  },
];

