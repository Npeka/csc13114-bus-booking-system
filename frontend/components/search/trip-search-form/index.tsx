"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";
import { format } from "date-fns";
import { LocationField } from "./location-field";
import { LocationPanel } from "./location-panel";
import { SwapLocationsButton } from "./swap-locations-button";
import { DatePickerField } from "./date-picker-field";
import { ReturnDatePickerField } from "./return-date-picker-field";
import { SharedDatePicker } from "./shared-date-picker";
import { PassengerField } from "./passenger-field";
import { PopularRoutes } from "./popular-routes";
import { VIETNAM_CITIES } from "./constants";
import { fuzzyMatchCity } from "./utils";
import { getCities } from "@/lib/api/constants-service";
import { formatDateForApi } from "@/lib/utils";

export function TripSearchForm() {
  const router = useRouter();
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  // Initialize as undefined to avoid hydration mismatch with PPR
  // Use function initializer to set date only on client side
  const [date, setDate] = useState<Date | undefined>(() => {
    if (typeof window !== "undefined") {
      return new Date();
    }
    return undefined;
  });
  const [isRoundTrip, setIsRoundTrip] = useState(false);
  const [returnDate, setReturnDate] = useState<Date | undefined>(undefined);
  const [passengers, setPassengers] = useState(1);
  const [recentLocations, setRecentLocations] = useState<string[]>([
    "Đà Lạt",
    "TP. Hồ Chí Minh",
  ]);
  const [locationPicker, setLocationPicker] = useState<{
    open: boolean;
    field: "origin" | "destination";
    search: string;
  }>({
    open: false,
    field: "origin",
    search: "",
  });
  const [datePicker, setDatePicker] = useState<{
    open: boolean;
    activeField: "departure" | "return";
  }>({
    open: false,
    activeField: "departure",
  });
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const dateFieldsRef = useRef<HTMLDivElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!date) return;

    const params = new URLSearchParams({
      from: origin,
      to: destination,
      date: formatDateForApi(date),
      passengers: passengers.toString(),
    });
    if (isRoundTrip && returnDate) {
      params.append("returnDate", formatDateForApi(returnDate));
    }
    router.push(`/trips?${params.toString()}`);
  };

  const handleToggleRoundTrip = () => {
    setIsRoundTrip((prev) => !prev);
    if (isRoundTrip) {
      // If turning off round trip, clear return date
      setReturnDate(undefined);
    } else {
      // If turning on, open the date picker immediately
      setDatePicker({ open: true, activeField: "return" });
      if (date) {
        // Set a default return date (e.g., 1 day after departure)
        const nextDay = new Date(date);
        nextDay.setDate(nextDay.getDate() + 1);
        setReturnDate(nextDay);
      }
    }
  };

  const openDatePicker = (field: "departure" | "return") => {
    // If clicking return date button (not in round trip mode), enable round trip first
    if (field === "return" && !isRoundTrip) {
      setIsRoundTrip(true);
      if (date) {
        const nextDay = new Date(date);
        nextDay.setDate(nextDay.getDate() + 1);
        setReturnDate(nextDay);
      }
    }
    setDatePicker({ open: true, activeField: field });
  };

  const closeDatePicker = () => {
    setDatePicker((prev) => ({ ...prev, open: false }));
  };

  const handleSwapLocations = () => {
    const temp = origin;
    setOrigin(destination);
    setDestination(temp);
  };

  const openLocationPicker = (field: "origin" | "destination") => {
    setLocationPicker({
      open: true,
      field,
      search: "",
    });
  };

  const closeLocationPicker = () => {
    setLocationPicker((prev) => ({ ...prev, open: false }));
  };

  const handleSelectLocation = (city: string) => {
    if (locationPicker.field === "origin") {
      setOrigin(city);
    } else {
      setDestination(city);
    }
    closeLocationPicker();
    setRecentLocations((prev) => {
      const next = [city, ...prev.filter((item) => item !== city)];
      return next.slice(0, 5);
    });
  };

  // Fetch cities from API (cached for 5 minutes)
  const { data: apiCities = [], isLoading: isLoadingCities } = useQuery({
    queryKey: ["cities"],
    queryFn: getCities,
    staleTime: 5 * 60 * 1000, // Cache for 5 minutes
  });

  const filteredLocations = useMemo(() => {
    const searchQuery = locationPicker.search.trim();
    
    // Determine the source: Use API cities if available, otherwise fallback to local VIETNAM_CITIES
    const citySource = apiCities.length > 0 ? apiCities : VIETNAM_CITIES;

    // If search query is empty, show all cities
    if (!searchQuery) {
      return citySource;
    }

    // Apply local fuzzy matching on the city source
    return citySource.filter((city) => fuzzyMatchCity(city, searchQuery));
  }, [locationPicker.search, apiCities]);

  useEffect(() => {
    if (locationPicker.open) {
      searchInputRef.current?.focus();
    }
  }, [locationPicker.open, locationPicker.field]);

  useEffect(() => {
    if (!locationPicker.open) {
      return;
    }

    const handleClick = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (
        target.closest("[data-location-trigger]") ||
        target.closest("[data-location-panel]")
      ) {
        return;
      }
      closeLocationPicker();
    };

    window.addEventListener("pointerdown", handleClick);
    return () => window.removeEventListener("pointerdown", handleClick);
  }, [locationPicker.open]);

  return (
    <Card className="w-full p-6 shadow-elevated md:p-8">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          {/* First Half: Origin and Destination */}
          <div className="relative grid grid-cols-1 gap-4 md:grid-cols-2">
            <LocationField
              id="origin"
              label="Điểm đi"
              placeholder="TP. Hồ Chí Minh"
              value={origin}
              onTrigger={() => openLocationPicker("origin")}
            >
              {locationPicker.open && locationPicker.field === "origin" && (
                <LocationPanel
                  searchInputRef={searchInputRef}
                  searchValue={locationPicker.search}
                  onSearchChange={(value) =>
                    setLocationPicker((prev) => ({
                      ...prev,
                      search: value,
                    }))
                  }
                  locations={filteredLocations}
                  onSelect={handleSelectLocation}
                  recentLocations={recentLocations}
                  isLoading={isLoadingCities}
                />
              )}
            </LocationField>

            <LocationField
              id="destination"
              label="Điểm đến"
              placeholder="Đà Nẵng"
              value={destination}
              iconClassName="text-primary"
              onTrigger={() => openLocationPicker("destination")}
            >
              {locationPicker.open &&
                locationPicker.field === "destination" && (
                  <LocationPanel
                    searchInputRef={searchInputRef}
                    searchValue={locationPicker.search}
                    onSearchChange={(value) =>
                      setLocationPicker((prev) => ({
                        ...prev,
                        search: value,
                      }))
                    }
                    locations={filteredLocations}
                    onSelect={handleSelectLocation}
                    recentLocations={recentLocations}
                    isLoading={isLoadingCities}
                  />
                )}
            </LocationField>

            <SwapLocationsButton onClick={handleSwapLocations} />
          </div>

          {/* Second Half: Date, Return Date, Passengers (3-3-2 layout) */}
          <div className="grid grid-cols-8 gap-4" ref={dateFieldsRef}>
            <div className="col-span-8 sm:col-span-3">
              <DatePickerField
                id="departure-date"
                label="Ngày đi"
                value={date}
                onClick={() => openDatePicker("departure")}
                isActive={
                  datePicker.open && datePicker.activeField === "departure"
                }
                required
              />
            </div>

            <div className="col-span-8 sm:col-span-3">
              <ReturnDatePickerField
                isRoundTrip={isRoundTrip}
                returnDate={returnDate}
                onClick={() => openDatePicker("return")}
                onToggle={handleToggleRoundTrip}
                isActive={
                  datePicker.open && datePicker.activeField === "return"
                }
              />
            </div>

            <div className="col-span-8 sm:col-span-2">
              <PassengerField
                value={passengers}
                onChange={(value) => setPassengers(value)}
              />
            </div>
          </div>
        </div>

        <Button
          type="submit"
          size="lg"
          className="h-12 w-full bg-primary text-base font-semibold text-white hover:bg-primary/90"
        >
          <Search className="mr-2 h-5 w-5" />
          Tìm chuyến xe
        </Button>
      </form>

      <PopularRoutes
        onSelectRoute={(route) => {
          setOrigin(route.from);
          setDestination(route.to);
        }}
      />

      {/* Shared Date Picker */}
      <SharedDatePicker
        isOpen={datePicker.open}
        onClose={closeDatePicker}
        departureDate={date}
        returnDate={returnDate}
        onDepartureDateChange={setDate}
        onReturnDateChange={setReturnDate}
        activeField={datePicker.activeField}
        triggerRef={dateFieldsRef}
      />
    </Card>
  );
}
