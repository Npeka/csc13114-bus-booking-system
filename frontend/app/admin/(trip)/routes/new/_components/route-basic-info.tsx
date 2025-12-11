"use client";

import { UseFormReturn } from "react-hook-form";
import { MapPin, Navigation, Clock } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface RouteFormData {
  origin: string;
  destination: string;
  distance_km: number;
  estimated_minutes: number;
}

interface RouteBasicInfoProps {
  form: UseFormReturn<RouteFormData>;
  cities: string[];
}

export function RouteBasicInfo({ form, cities }: RouteBasicInfoProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin tuyến đường</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Origin and Destination - Same Row */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <FormField
            control={form.control}
            name="origin"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <MapPin className="mr-2 inline h-4 w-4" />
                  Điểm đi
                </FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Chọn điểm đi" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {cities.map((city) => (
                      <SelectItem key={city} value={city}>
                        {city}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="destination"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Navigation className="mr-2 inline h-4 w-4" />
                  Điểm đến
                </FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Chọn điểm đến" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {cities.map((city) => (
                      <SelectItem key={city} value={city}>
                        {city}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        {/* Distance and Duration - Same Row */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <FormField
            control={form.control}
            name="distance_km"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Navigation className="mr-2 inline h-4 w-4" />
                  Khoảng cách (km)
                </FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    step="0.1"
                    min="0.1"
                    placeholder="VD: 308"
                    {...field}
                    onChange={(e) =>
                      field.onChange(parseFloat(e.target.value) || 0)
                    }
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="estimated_minutes"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <Clock className="mr-2 inline h-4 w-4" />
                  Thời gian ước tính (phút)
                </FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    min="1"
                    placeholder="VD: 420"
                    {...field}
                    onChange={(e) =>
                      field.onChange(parseInt(e.target.value) || 0)
                    }
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
      </CardContent>
    </Card>
  );
}
