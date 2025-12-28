"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Star, Clock, Users, ChevronDown, ChevronUp } from "lucide-react";
import { getTripReviews, getTripReviewSummary } from "@/lib/api/booking";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

export interface Trip {
  id: string;
  operator: string;
  operatorRating: number;
  operatorLogo?: string;
  departureTime: string;
  arrivalTime: string;
  duration: string;
  origin: string;
  destination: string;
  price: number;
  availableSeats: number;
  busType: string;
  amenities: string[];
}

interface TripCardProps {
  trip: Trip;
  onSelect: (tripId: string) => void;
}

export function TripCard({ trip, onSelect }: TripCardProps) {
  const [showReviews, setShowReviews] = useState(false);

  // Fetch review summary
  const { data: reviewSummary } = useQuery({
    queryKey: ["trip-review-summary", trip.id],
    queryFn: () => getTripReviewSummary(trip.id),
    enabled: showReviews,
  });

  // Fetch reviews when expanded
  const { data: reviewsData, isLoading: reviewsLoading } = useQuery({
    queryKey: ["trip-reviews", trip.id],
    queryFn: () => getTripReviews(trip.id, 1, 5),
    enabled: showReviews,
  });

  return (
    <Card className="card-hover cursor-pointer py-0!">
      <CardContent className="p-4" onClick={() => onSelect(trip.id)}>
        {/* Top Row: Operator Info + Price/Action */}
        <div className="mb-2 flex items-center justify-between gap-4">
          {/* Operator Info */}
          <div className="flex min-w-0 flex-1 items-center space-x-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-secondary">
              <span className="text-lg">🚌</span>
            </div>
            <div className="min-w-0 flex-1">
              <h3
                className="truncate text-base font-semibold"
                title={trip.operator}
              >
                {trip.operator}
              </h3>
              <div className="flex items-center text-xs text-muted-foreground">
                <Star className="mr-1 h-3 w-3 shrink-0 fill-warning text-warning" />
                <span>{trip.operatorRating.toFixed(1)}</span>
                <span className="mx-1">•</span>
                <span className="truncate">{trip.busType}</span>
              </div>
            </div>
          </div>

          {/* Price and Action */}
          <div className="flex shrink-0 items-center gap-4">
            <div className="text-right">
              <div className="text-xs text-muted-foreground">Giá từ</div>
              <div className="text-xl font-bold text-primary">
                {trip.price.toLocaleString()}đ
              </div>
              <div className="flex items-center justify-end text-xs text-muted-foreground">
                <Users className="mr-1 h-3 w-3" />
                <span>{trip.availableSeats} chỗ</span>
              </div>
            </div>
            <Button
              className="bg-primary text-white hover:bg-primary/90"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onSelect(trip.id);
              }}
            >
              Chọn chỗ
            </Button>
          </div>
        </div>

        {/* Bottom Row: Trip Time Details */}
        <div className="flex items-center justify-between gap-6 px-2 py-2">
          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-primary">
              {trip.departureTime}
            </div>
            <div className="mt-1 text-xs font-medium text-muted-foreground">
              {trip.origin}
            </div>
          </div>

          <div className="flex flex-col items-center justify-center px-5 py-1.5">
            <div className="mb-1 flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-primary"></div>
              <div className="h-0.5 w-12 bg-linear-to-r from-primary to-primary/70"></div>
              <Clock className="h-4 w-4 text-primary" />
              <div className="h-0.5 w-12 bg-linear-to-r from-primary/70 to-primary"></div>
              <div className="h-2 w-2 rounded-full bg-primary"></div>
            </div>
            <div className="text-xs font-semibold whitespace-nowrap text-primary">
              {trip.duration}
            </div>
          </div>

          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-primary">
              {trip.arrivalTime}
            </div>
            <div className="mt-1 text-xs font-medium text-muted-foreground">
              {trip.destination}
            </div>
          </div>
        </div>

        {/* Amenities */}
        {trip.amenities.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-2 pt-2">
            {trip.amenities.map((amenity, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                {amenity}
              </Badge>
            ))}
          </div>
        )}

        {/* Reviews Toggle */}
        <div className="mt-3 border-t pt-3">
          <button
            className="flex w-full items-center justify-between text-sm font-medium hover:text-primary"
            onClick={(e) => {
              e.stopPropagation();
              setShowReviews(!showReviews);
            }}
          >
            <span className="flex items-center gap-2">
              <Star className="h-4 w-4 fill-warning text-warning" />
              <span>
                {reviewSummary
                  ? `${reviewSummary.average_rating.toFixed(1)} (${reviewSummary.total_reviews} đánh giá)`
                  : "Xem đánh giá"}
              </span>
            </span>
            {showReviews ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </button>

          {/* Reviews Content */}
          {showReviews && (
            <div className="mt-3 space-y-3">
              {reviewsLoading ? (
                <div className="space-y-2">
                  <Skeleton className="h-16 w-full" />
                  <Skeleton className="h-16 w-full" />
                </div>
              ) : reviewsData?.data && reviewsData.data.length > 0 ? (
                <>
                  {/* Rating Distribution */}
                  {reviewSummary && (
                    <div className="rounded-lg bg-secondary/50 p-3">
                      <div className="grid grid-cols-5 gap-2 text-center text-xs">
                        {[5, 4, 3, 2, 1].map((rating) => {
                          const count = reviewSummary[
                            `rating_${rating}_count` as keyof typeof reviewSummary
                          ] as number;
                          return (
                            <div key={rating}>
                              <div className="flex items-center justify-center gap-0.5">
                                {rating}
                                <Star className="h-3 w-3 fill-warning text-warning" />
                              </div>
                              <div className="mt-1 font-semibold">{count}</div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  )}

                  {/* Reviews List */}
                  <div className="space-y-3">
                    {reviewsData.data.slice(0, 3).map((review) => (
                      <div
                        key={review.id}
                        className="rounded-lg border bg-card p-3"
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex items-center gap-1">
                            {Array.from({ length: 5 }).map((_, i) => (
                              <Star
                                key={i}
                                className={`h-3 w-3 ${
                                  i < review.rating
                                    ? "fill-warning text-warning"
                                    : "text-muted-foreground"
                                }`}
                              />
                            ))}
                          </div>
                          <span className="text-xs text-muted-foreground">
                            {format(new Date(review.created_at), "dd/MM/yyyy", {
                              locale: vi,
                            })}
                          </span>
                        </div>
                        {review.comment && (
                          <p className="mt-2 text-sm text-muted-foreground">
                            {review.comment}
                          </p>
                        )}
                        {review.is_verified && (
                          <Badge variant="outline" className="mt-2 text-xs">
                            <span className="mr-1">✓</span>
                            Đã xác thực
                          </Badge>
                        )}
                      </div>
                    ))}
                  </div>

                  {/* View More Link */}
                  {reviewsData.meta.total > 3 && (
                    <button
                      className="text-sm text-primary hover:underline"
                      onClick={(e) => {
                        e.stopPropagation();
                        onSelect(trip.id);
                      }}
                    >
                      Xem tất cả {reviewsData.meta.total} đánh giá →
                    </button>
                  )}
                </>
              ) : (
                <p className="text-center text-sm text-muted-foreground">
                  Chưa có đánh giá
                </p>
              )}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export function TripCardSkeleton() {
  return (
    <Card>
      <CardContent className="p-3">
        {/* Top Row */}
        <div className="mb-2 flex items-center justify-between gap-4">
          <div className="flex flex-1 items-center space-x-3">
            <Skeleton className="h-10 w-10 shrink-0 rounded-lg" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-4 w-32" />
              <Skeleton className="h-3 w-24" />
            </div>
          </div>
          <div className="flex shrink-0 items-center gap-4">
            <div className="space-y-2">
              <Skeleton className="ml-auto h-3 w-12" />
              <Skeleton className="h-5 w-20" />
              <Skeleton className="ml-auto h-3 w-16" />
            </div>
            <Skeleton className="h-9 w-20" />
          </div>
        </div>
        {/* Bottom Row */}
        <div className="flex items-center justify-between gap-6 px-2 py-2">
          <Skeleton className="mx-auto h-12 w-16" />
          <Skeleton className="h-10 w-32 rounded-lg" />
          <Skeleton className="mx-auto h-12 w-16" />
        </div>
      </CardContent>
    </Card>
  );
}
