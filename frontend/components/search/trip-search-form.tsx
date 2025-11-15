"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card } from "@/components/ui/card";
import { Calendar, MapPin, Users, ArrowRightLeft, Search } from "lucide-react";
import { format } from "date-fns";

export function TripSearchForm() {
  const router = useRouter();
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const [date, setDate] = useState(format(new Date(), "yyyy-MM-dd"));
  const [passengers, setPassengers] = useState(1);

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

  return (
    <Card className="w-full max-w-5xl p-6 shadow-elevated md:p-8">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
          {/* Origin */}
          <div className="relative space-y-2">
            <Label htmlFor="origin" className="text-sm font-semibold">
              Điểm đi
            </Label>
            <div className="relative">
              <MapPin className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
              <Input
                id="origin"
                type="text"
                placeholder="TP. Hồ Chí Minh"
                value={origin}
                onChange={(e) => setOrigin(e.target.value)}
                className="h-12 pl-10"
                required
              />
            </div>
          </div>

          {/* Swap Button - Hidden on mobile, visible on desktop */}
          <div className="hidden items-end lg:flex">
            <Button
              type="button"
              variant="ghost"
              size="icon"
              onClick={handleSwapLocations}
              className="h-12 w-12"
              aria-label="Đổi điểm đi đến"
            >
              <ArrowRightLeft className="h-5 w-5" />
            </Button>
          </div>

          {/* Destination */}
          <div className="space-y-2 lg:col-start-3">
            <Label htmlFor="destination" className="text-sm font-semibold">
              Điểm đến
            </Label>
            <div className="relative">
              <MapPin className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-brand-primary" />
              <Input
                id="destination"
                type="text"
                placeholder="Đà Nẵng"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                className="h-12 pl-10"
                required
              />
            </div>
          </div>

          {/* Date */}
          <div className="space-y-2">
            <Label htmlFor="date" className="text-sm font-semibold">
              Ngày đi
            </Label>
            <div className="relative">
              <Calendar className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
              <Input
                id="date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="h-12 pl-10"
                min={format(new Date(), "yyyy-MM-dd")}
                required
              />
            </div>
          </div>

          {/* Passengers */}
          <div className="space-y-2">
            <Label htmlFor="passengers" className="text-sm font-semibold">
              Số hành khách
            </Label>
            <div className="relative">
              <Users className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
              <Input
                id="passengers"
                type="number"
                min="1"
                max="10"
                value={passengers}
                onChange={(e) => setPassengers(parseInt(e.target.value) || 1)}
                className="h-12 pl-10"
                required
              />
            </div>
          </div>
        </div>

        {/* Search Button */}
        <Button
          type="submit"
          size="lg"
          className="w-full bg-brand-primary text-white hover:bg-brand-primary-hover h-12 text-base font-semibold"
        >
          <Search className="mr-2 h-5 w-5" />
          Tìm chuyến xe
        </Button>
      </form>

      {/* Popular Routes */}
      <div className="mt-6 border-t pt-6">
        <p className="mb-3 text-sm font-medium text-muted-foreground">
          Tuyến đường phổ biến:
        </p>
        <div className="flex flex-wrap gap-2">
          {popularRoutes.map((route) => (
            <Button
              key={route.id}
              type="button"
              variant="outline"
              size="sm"
              onClick={() => {
                setOrigin(route.from);
                setDestination(route.to);
              }}
              className="text-xs"
            >
              {route.from} → {route.to}
            </Button>
          ))}
        </div>
      </div>
    </Card>
  );
}

const popularRoutes = [
  { id: 1, from: "Hà Nội", to: "Đà Nẵng" },
  { id: 2, from: "TP. Hồ Chí Minh", to: "Đà Lạt" },
  { id: 3, from: "Hà Nội", to: "Sa Pa" },
  { id: 4, from: "TP. Hồ Chí Minh", to: "Nha Trang" },
  { id: 5, from: "Hà Nội", to: "Hạ Long" },
];

