"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { StarRatingInput } from "./star-rating-input";
import { createReview, CreateReviewRequest } from "@/lib/api/booking";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

interface ReviewFormProps {
  bookingId: string;
  tripId: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function ReviewForm({
  bookingId,
  tripId,
  onSuccess,
  onCancel,
}: ReviewFormProps) {
  const [rating, setRating] = useState(0);
  const [comment, setComment] = useState("");
  const queryClient = useQueryClient();

  const { mutate: submitReview, isPending } = useMutation({
    mutationFn: (data: CreateReviewRequest) => createReview(bookingId, data),
    onSuccess: () => {
      toast.success("Đánh giá đã được gửi thành công!");
      queryClient.invalidateQueries({
        queryKey: ["booking-review", bookingId],
      });
      queryClient.invalidateQueries({ queryKey: ["trip-reviews", tripId] });
      queryClient.invalidateQueries({
        queryKey: ["trip-review-summary", tripId],
      });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể gửi đánh giá. Vui lòng thử lại.");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (rating === 0) {
      toast.error("Vui lòng chọn số sao đánh giá");
      return;
    }
    submitReview({
      booking_id: bookingId,
      rating,
      comment: comment.trim() || undefined,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label>Đánh giá của bạn *</Label>
        <StarRatingInput
          value={rating}
          onChange={setRating}
          size="lg"
          disabled={isPending}
        />
        {rating > 0 && (
          <p className="text-sm text-muted-foreground">
            {rating === 5 && "Tuyệt vời!"}
            {rating === 4 && "Rất tốt"}
            {rating === 3 && "Bình thường"}
            {rating === 2 && "Không hài lòng"}
            {rating === 1 && "Rất tệ"}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="comment">Nhận xét (tùy chọn)</Label>
        <Textarea
          id="comment"
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          placeholder="Chia sẻ trải nghiệm của bạn về chuyến đi..."
          maxLength={1000}
          rows={4}
          disabled={isPending}
        />
        <p className="text-xs text-muted-foreground">
          {comment.length}/1000 ký tự
        </p>
      </div>

      <div className="flex gap-2">
        {onCancel && (
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isPending}
          >
            Hủy
          </Button>
        )}
        <Button type="submit" disabled={isPending || rating === 0}>
          {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Gửi đánh giá
        </Button>
      </div>
    </form>
  );
}
