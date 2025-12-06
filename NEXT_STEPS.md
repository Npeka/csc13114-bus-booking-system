# Next Steps - Week 4 Implementation

## 📋 Overview

**Focus:** Payment integration, notifications, and post-booking management to complete the transaction lifecycle.

**Estimated Duration:** 29-36 hours  
**Target Completion:** End of Week 4

---

## ✅ Prerequisites (Week 3 Completed)

Before starting Week 4, ensure the following are complete:

- ✅ Interactive seat map component
- ✅ Seat locking mechanism (temporary reservations)
- ✅ Real-time seat availability updates
- ✅ Booking creation and management APIs
- ✅ Passenger information collection forms
- ✅ Booking history and management dashboard
- ✅ Guest checkout flow
- ✅ E-ticket generation and delivery

---

## 🎯 Week 4 Tasks

### 1. User Portal / Payments (~8 hours)

#### Integrate PayOS Payment Gateway (~3 hours)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create PayOS account and obtain API credentials
- [ ] Implement PayOS API client in `payment-service`
- [ ] Create payment initiation endpoint `POST /api/v1/payments/create`
- [ ] Handle credit card and digital wallet payment methods
- [ ] Test payment flow in sandbox environment

**Files to create/modify:**

- `backend/payment-service/internal/service/payos_service.go`
- `backend/payment-service/internal/handler/payment_handler.go`
- `frontend/lib/api/payment-service.ts`

---

#### Implement Payment Webhook Handling (~3 hours)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create webhook endpoint `POST /api/v1/payments/webhook`
- [ ] Verify webhook signatures for security
- [ ] Update booking status based on payment result
- [ ] Handle payment success, failure, and pending states
- [ ] Implement idempotency for webhook processing

**Files to create/modify:**

- `backend/payment-service/internal/handler/webhook_handler.go`
- `backend/payment-service/internal/service/webhook_service.go`

---

#### Create Payment Confirmation and Failure Flows (~2 hours)

**Assigned to:** Frontend developer

**Tasks:**

- [ ] Create payment success page `/payment/success`
- [ ] Create payment failure page `/payment/failure`
- [ ] Implement automatic redirect to booking details on success
- [ ] Add retry payment button on failure page
- [ ] Show payment receipt with transaction details

**Files to create:**

- `frontend/app/payment/success/page.tsx`
- `frontend/app/payment/failure/page.tsx`
- `frontend/components/payment/payment-receipt.tsx`

---

### 2. User Portal / Notifications (~9 hours)

#### Setup Email Service (~1 hour)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create SendGrid or AWS SES account
- [ ] Configure SMTP credentials in environment variables
- [ ] Create email service wrapper in `shared/email`
- [ ] Test email sending functionality

**Files to create:**

- `backend/shared/email/service.go`
- `backend/shared/email/templates.go`

---

#### Create Email Templates for Booking Confirmations (~2 hours)

**Assigned to:** Frontend developer

**Tasks:**

- [ ] Design HTML email template for booking confirmation
- [ ] Create template for payment receipt
- [ ] Create template for booking cancellation
- [ ] Add company branding and essential trip information
- [ ] Make templates mobile-responsive

**Files to create:**

- `backend/shared/email/templates/booking_confirmation.html`
- `backend/shared/email/templates/payment_receipt.html`
- `backend/shared/email/templates/booking_cancellation.html`

---

#### Implement SMS Notifications (Optional) (~2 hours)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create Twilio account (optional)
- [ ] Implement SMS service wrapper
- [ ] Create SMS template for booking confirmation
- [ ] Add SMS delivery for critical notifications
- [ ] Handle SMS delivery failures gracefully

**Files to create:**

- `backend/shared/sms/service.go` (optional)

---

#### Setup Trip Reminder Notifications (~2 hours)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create scheduled job using cron or similar
- [ ] Query bookings with trips departing in 24 hours
- [ ] Send reminder emails to passengers
- [ ] Send SMS reminders (if implemented)
- [ ] Log notification delivery status

**Files to create/modify:**

- `backend/booking-service/internal/jobs/reminder_job.go`
- `backend/booking-service/cmd/worker/main.go` (if using separate worker)

---

#### Create Notification Preferences Management (~2 hours)

**Assigned to:** Frontend developer

**Tasks:**

- [ ] Add notification preferences to user profile
- [ ] Create UI to toggle email/SMS notifications
- [ ] Add API endpoints to update preferences
- [ ] Save preferences in user database

**Files to create/modify:**

- `frontend/app/(auth)/profile/page.tsx` (add preferences section)
- `backend/user-service/internal/model/user.go` (add notification fields)
- `backend/user-service/internal/handler/user_handler.go`

---

### 3. User Portal / Management (~5 hours)

#### Create Booking Modification Functionality (~4 hours)

**Assigned to:** Full-stack developer

**Tasks:**

- [ ] Add modify booking button to booking details
- [ ] Allow passenger detail changes
- [ ] Implement seat change if available
- [ ] Validate modification constraints (time limits)
- [ ] Update booking and send confirmation

**Files to create/modify:**

- `frontend/app/(auth)/bookings/[id]/edit/page.tsx`
- `backend/booking-service/internal/handler/booking_handler.go` (add update endpoint)

---

#### Setup Automated Booking Expiration (~1 hour)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create scheduled job to check expired bookings
- [ ] Cancel bookings unpaid after timeout (e.g., 15 minutes)
- [ ] Release seats back to availability
- [ ] Send cancellation notification
- [ ] Log expiration events

**Files to create:**

- `backend/booking-service/internal/jobs/expiration_job.go`

---

### 4. Admin Portal (~5 hours)

#### Create Revenue Analytics Dashboard (~3 hours)

**Assigned to:** Frontend developer

**Tasks:**

- [ ] Create revenue dashboard UI in admin portal
- [ ] Add charts for daily/weekly/monthly revenue
- [ ] Show revenue by route, operator, bus type
- [ ] Add total bookings and cancellation stats
- [ ] Implement date range filters

**Files to create:**

- `frontend/app/admin/analytics/revenue/page.tsx`
- `frontend/components/admin/revenue-chart.tsx`

---

#### Implement Booking Analytics and Reporting (~2 hours)

**Assigned to:** Backend developer

**Tasks:**

- [ ] Create analytics endpoint `GET /api/v1/admin/analytics`
- [ ] Aggregate booking data by time periods
- [ ] Calculate conversion rates
- [ ] Identify popular routes and trends
- [ ] Return data for dashboard charts

**Files to create:**

- `backend/booking-service/internal/handler/analytics_handler.go`
- `backend/booking-service/internal/service/analytics_service.go`

---

### 5. System & Infrastructure (~2 hours)

#### Setup Real-time Monitoring Dashboard (~2 hours)

**Assigned to:** DevOps engineer

**Tasks:**

- [ ] Configure Grafana or similar monitoring tool
- [ ] Create dashboard for key metrics:
  - API response times
  - Error rates
  - Database connection pool
  - Payment success rate
  - Email delivery rate
- [ ] Setup alerts for critical thresholds
- [ ] Document monitoring setup

**Files/Tools:**

- Grafana dashboards (JSON config)
- Prometheus metrics endpoints
- AlertManager rules

---

## 📊 Progress Tracking

| Category       | Tasks  | Estimated Hours | Assigned           | Status         |
| -------------- | ------ | --------------- | ------------------ | -------------- |
| Payments       | 3      | 8h              | Backend + Frontend | ⬜ Not Started |
| Notifications  | 5      | 9h              | Backend + Frontend | ⬜ Not Started |
| Management     | 2      | 5h              | Full-stack         | ⬜ Not Started |
| Admin Portal   | 2      | 5h              | Backend + Frontend | ⬜ Not Started |
| Infrastructure | 1      | 2h              | DevOps             | ⬜ Not Started |
| **Total**      | **13** | **29h**         |                    |                |

---

## 🔑 External Dependencies

### Required Services

- **PayOS** - Payment gateway account and API credentials
- **SendGrid/AWS SES** - Email delivery service
- **Twilio** (Optional) - SMS notifications
- **Grafana/Prometheus** - Monitoring and observability

### Environment Variables to Add

```bash
# Payment Service
PAYOS_API_KEY=your_api_key
PAYOS_MERCHANT_ID=your_merchant_id
PAYOS_WEBHOOK_SECRET=your_webhook_secret

# Email Service
SENDGRID_API_KEY=your_sendgrid_key
# or
AWS_SES_ACCESS_KEY=your_aws_key
AWS_SES_SECRET_KEY=your_aws_secret
EMAIL_FROM=noreply@busticket.vn

# SMS Service (Optional)
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=your_twilio_number
```

---

## 🎓 Technical Guidelines

### Payment Integration

- Use existing `apiClient` pattern for PayOS API calls
- Follow `booking-service.ts` structure for API client
- Maintain consistent error handling with `handleApiError`
- Add proper TypeScript types for all payment responses
- Implement retry logic for failed webhook deliveries

### Notification System

- Use template engine for dynamic email content
- Store email/SMS templates in database or files
- Implement queue system for asynchronous delivery
- Log all notification attempts for debugging
- Handle bounce/failure tracking

### Analytics

- Use database views or CTEs for complex aggregations
- Cache dashboard data for better performance
- Implement pagination for large datasets
- Export functionality for reports (CSV/PDF)

### Security

- Validate all webhook signatures
- Use HTTPS for all payment communications
- Never log sensitive payment data
- Implement rate limiting on payment endpoints
- Sanitize all user inputs in email templates

---

## 🚀 Getting Started

1. **Setup External Services**

   - Create PayOS merchant account
   - Configure email service provider
   - (Optional) Setup SMS provider

2. **Update Environment Variables**

   - Add credentials to `.env.dev` files
   - Update Kubernetes secrets for production

3. **Start with Payments**

   - Begin with backend PayOS integration
   - Test in sandbox mode thoroughly
   - Create frontend payment flow
   - Implement webhook handling

4. **Implement Notifications**

   - Setup email service
   - Create templates
   - Test email delivery
   - Add to booking confirmation flow

5. **Build Management Features**

   - Add booking modification UI
   - Implement expiration job
   - Test edge cases

6. **Create Admin Analytics**

   - Build backend aggregation endpoints
   - Create dashboard visualizations
   - Add export functionality

7. **Setup Monitoring**
   - Configure metrics collection
   - Create dashboards
   - Set up alerts

---

## ✅ Definition of Done

Week 4 is considered complete when:

- [ ] Users can pay for bookings via PayOS
- [ ] Payment webhooks update booking status correctly
- [ ] Confirmation emails sent automatically
- [ ] Users can modify bookings within allowed time
- [ ] Unpaid bookings expire and release seats
- [ ] Admin dashboard shows revenue analytics
- [ ] Monitoring dashboard tracks system health
- [ ] All APIs have proper error handling
- [ ] Payment flow tested end-to-end
- [ ] Documentation updated

---

## 📞 Support

For questions or issues during Week 4 implementation:

- Review PayOS documentation: https://payos.vn/docs
- Check SendGrid API docs: https://docs.sendgrid.com
- Refer to existing code patterns in Week 2-3 implementations
- Ask team for code review and pair programming help

**Good luck with Week 4! 🚀**
