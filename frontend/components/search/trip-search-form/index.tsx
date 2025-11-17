"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";
import { format } from "date-fns";
import { LocationField } from "./location-field";
import { LocationPanel } from "./location-panel";
import { SwapLocationsButton } from "./swap-locations-button";
import { DateField } from "./date-field";
import { PassengerField } from "./passenger-field";
import { PopularRoutes } from "./popular-routes";
import { VIETNAM_CITIES } from "./constants";
import { fuzzyMatchCity } from "./utils";

export function TripSearchForm() {
  const router = useRouter();
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const [date, setDate] = useState(format(new Date(), "yyyy-MM-dd"));
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
  const searchInputRef = useRef<HTMLInputElement | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const params = new URLSearchParams({
      from: origin,
      to: destination,
      date: date,
      passengers: passengers.toString(),
    });
    router.push(`/trips?${params.toString()}`);
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

  const filteredLocations = useMemo(() => {
    if (!locationPicker.search.trim()) {
      return VIETNAM_CITIES;
    }
    return VIETNAM_CITIES.filter((city) =>
      fuzzyMatchCity(city, locationPicker.search),
    );
  }, [locationPicker.search]);

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
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
          <div className="md:col-span-2 lg:col-span-2">
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
                  />
                )}
              </LocationField>

              <LocationField
                id="destination"
                label="Điểm đến"
                placeholder="Đà Nẵng"
                value={destination}
                iconClassName="text-brand-primary"
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
                    />
                  )}
              </LocationField>

              <SwapLocationsButton onClick={handleSwapLocations} />
            </div>
          </div>

          <DateField value={date} onChange={setDate} />

          <PassengerField
            value={passengers}
            onChange={(value) => setPassengers(value)}
          />
        </div>

        <Button
          type="submit"
          size="lg"
          className="h-12 w-full bg-brand-primary text-base font-semibold text-white hover:bg-brand-primary-hover"
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
    </Card>
  );
}
