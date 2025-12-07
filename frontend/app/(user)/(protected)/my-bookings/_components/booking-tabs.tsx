import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Download, X } from "lucide-react";
import Link from "next/link";
import { BookingCard } from "./booking-card";

type UIBooking = {
  id: string;
  bookingReference: string;
  status: string;
  trip: {
    operator: string;
    origin: string;
    destination: string;
    date: string;
    departureTime: string;
  };
  seats: string[];
  price: number;
  refundAmount?: number;
  cancelledReason?: string;
};

interface BookingTabsProps {
  upcomingBookings: UIBooking[];
  pastBookings: UIBooking[];
  cancelledBookings: UIBooking[];
}

export function BookingTabs({
  upcomingBookings,
  pastBookings,
  cancelledBookings,
}: BookingTabsProps) {
  return (
    <Tabs defaultValue="upcoming" className="space-y-4">
      <TabsList>
        <TabsTrigger value="upcoming">
          Sắp diễn ra ({upcomingBookings.length})
        </TabsTrigger>
        <TabsTrigger value="past">
          Đã hoàn thành ({pastBookings.length})
        </TabsTrigger>
        <TabsTrigger value="cancelled">
          Đã hủy ({cancelledBookings.length})
        </TabsTrigger>
      </TabsList>

      {/* Upcoming Bookings */}
      <TabsContent value="upcoming" className="space-y-4">
        {upcomingBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Bạn chưa có chuyến đi nào sắp tới
              </p>
              <Button asChild className="mt-4">
                <Link href="/">Đặt vé ngay</Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {upcomingBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <>
                    <Button variant="outline" size="sm">
                      <Download className="h-4 w-4" />
                      Tải vé
                    </Button>
                    <Button variant="outline" size="sm">
                      <X className="h-4 w-4" />
                      Hủy vé
                    </Button>
                  </>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>

      {/* Past Bookings */}
      <TabsContent value="past" className="space-y-4">
        {pastBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Chưa có chuyến đi nào đã hoàn thành
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {pastBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <>
                    <Button variant="outline" size="sm">
                      <Download className="mr-2 h-4 w-4" />
                      Tải vé
                    </Button>
                    <Button variant="outline" size="sm">
                      Đặt lại
                    </Button>
                  </>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>

      {/* Cancelled Bookings */}
      <TabsContent value="cancelled" className="space-y-4">
        {cancelledBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Chưa có vé nào bị hủy</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {cancelledBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <Button variant="outline" size="sm">
                    Đặt lại
                  </Button>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>
    </Tabs>
  );
}
