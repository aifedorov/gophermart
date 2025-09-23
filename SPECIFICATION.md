# Technical Specification

## Gophermart Loyalty System

---

### General Requirements

The system is an HTTP API with the following business logic requirements:

* user registration, authentication, and authorization;
* accepting order numbers from registered users;
* tracking and maintaining a list of order numbers submitted by registered users;
* tracking and maintaining loyalty point balances for registered users;
* verifying submitted order numbers through the loyalty points calculation system;
* crediting appropriate rewards to user loyalty accounts for each valid order number.

![image](https://pictures.s3.yandex.net:443/resources/gophermart2x_1634502166.png)

### Abstract System Interaction Flow

The following abstract business logic describes user interaction with the system:

1. User registers in the Gophermart loyalty system.
2. User makes a purchase in the Gophermart online store.
3. Order is submitted to the loyalty points calculation system.
4. User submits the completed order number to the loyalty system.
5. System links the order number with the user and verifies the number with the loyalty points calculation system.
6. If a positive loyalty points calculation exists, loyalty points are credited to the user's account.
7. User withdraws available loyalty points for partial or full payment of subsequent orders in the Gophermart online store.

Notes:

- Step 2 is hypothetical and does not require implementation in this work;
- Step 3 is implemented in the loyalty points calculation system and does not require implementation in this work.

### Loyalty Points Calculation System

The loyalty points calculation system is an external service within a trusted perimeter. It operates as a black box and is not accessible for inspection by external clients. The system calculates due loyalty points for completed orders using complex algorithms that may change at any time.

External consumers only have access to information about the number of loyalty points due for a specific order. The reasons for the presence or absence of accruals are unknown to external consumers.

The interaction protocol with the service will be provided at the end.

### HTTP API Summary

The Gophermart loyalty system must provide the following HTTP handlers:

* `POST /api/user/register` — user registration;
* `POST /api/user/login` — user authentication;
* `POST /api/user/orders` — user order number submission for calculation;
* `GET /api/user/orders` — retrieving a list of order numbers submitted by the user, their processing statuses, and accrual information;
* `GET /api/user/balance` — retrieving the current balance of the user's loyalty account;
* `POST /api/user/balance/withdraw` — request to withdraw points from the loyalty account for payment of a new order;
* `GET /api/user/withdrawals` — retrieving information about withdrawals from the user's loyalty account.

### General Constraints and Requirements

* data storage — PostgreSQL;
* table structure is at the student's discretion;
* data types and storage formats (including passwords and other sensitive information) are at the student's discretion;
* client may support HTTP requests/responses with data compression;
* client is not obligated to make requests according to the API specification below, any request validation is at the student's discretion;
* format and algorithm for user authentication and authorization verification are at the student's discretion;
* order numbers are unique and never repeat;
* an order number can only be accepted for processing once from one user;
* an order number may not have any accrual;
* rewards are credited and spent in virtual points at a rate of 1 point = 1 ruble.

#### **User Registration**

Handler: `POST /api/user/register`.

Registration is performed using a login/password pair. Each login must be unique.
After successful registration, automatic user authentication should occur.

Request format:

```
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Possible response codes:

- `200` — user successfully registered and authenticated;
- `400` — invalid request format;
- `409` — login already taken;
- `500` — internal server error.

#### **User Authentication**

Handler: `POST /api/user/login`.

Authentication is performed using a login/password pair.

Request format:

```
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Possible response codes:

- `200` — user successfully authenticated;
- `400` — invalid request format;
- `401` — invalid login/password pair;
- `500` — internal server error.

#### **Order Number Submission**

Handler: `POST /api/user/orders`.

Handler is available only to authenticated users. An order number is a sequence of digits of arbitrary length.

Order number can be validated for input correctness using the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm){target="_blank"}.

Request format:

```
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903
```

Possible response codes:

- `200` — order number was already uploaded by this user;
- `202` — new order number accepted for processing;
- `400` — invalid request format;
- `401` — user not authenticated;
- `409` — order number was already uploaded by another user;
- `422` — invalid order number format;
- `500` — internal server error.

#### **Retrieving List of Submitted Order Numbers**

Handler: `GET /api/user/orders`.

Handler is available only to authorized users. Order numbers in the response should be sorted by upload time from newest to oldest. Date format — RFC3339.

Available calculation processing statuses:

- `NEW` — order uploaded to system but not yet in processing;
- `PROCESSING` — order reward is being calculated;
- `INVALID` — reward calculation system rejected the calculation;
- `PROCESSED` — order data verified and calculation information successfully obtained.

Request format:

```
GET /api/user/orders HTTP/1.1
Content-Length: 0
```

Possible response codes:

- `200` — successful request processing.

  Response format:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
    	{
            "number": "9278923470",
            "status": "PROCESSED",
            "accrual": 500,
            "uploaded_at": "2020-12-10T15:15:45+03:00"
        },
        {
            "number": "12345678903",
            "status": "PROCESSING",
            "uploaded_at": "2020-12-10T15:12:01+03:00"
        },
        {
            "number": "346436439",
            "status": "INVALID",
            "uploaded_at": "2020-12-09T16:09:53+03:00"
        }
    ]
    ```

- `204` — no data for response.
- `401` — user not authorized.
- `500` — internal server error.

#### **Retrieving Current User Balance**

Handler: `GET /api/user/balance`.

Handler is available only to authorized users. The response should contain data about the current amount of loyalty points, as well as the amount of points used during the entire registration period.

Request format:

```
GET /api/user/balance HTTP/1.1
Content-Length: 0
```

Possible response codes:

- `200` — successful request processing.

  Response format:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
    	"current": 500.5,
    	"withdrawn": 42
    }
    ```

- `401` — user not authorized.
- `500` — internal server error.

#### **Withdrawal Request**

Handler: `POST /api/user/balance/withdraw`

Handler is available only to authorized users. The order number represents a hypothetical number of a new user order for which points are being withdrawn.

Note: for successful withdrawal, successful request registration is sufficient; no external accrual systems are provided and do not need to be implemented.

Request format:

```
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
	"order": "2377225624",
    "sum": 751
}
```

Here `order` is the order number, and `sum` is the amount of points to withdraw for payment.

Possible response codes:

- `200` — successful request processing;
- `401` — user not authorized;
- `402` — insufficient funds;
- `422` — invalid order number;
- `500` — internal server error.

#### **Retrieving Withdrawal Information**

Handler: `GET /api/user/withdrawals`.

Handler is available only to authorized users. Withdrawal records in the response should be sorted by withdrawal time from newest to oldest. Date format — RFC3339.

Request format:

```
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
```

Possible response codes:

- `200` — successful request processing.

  Response format:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
        {
            "order": "2377225624",
            "sum": 500,
            "processed_at": "2020-12-09T16:09:57+03:00"
        }
    ]
    ```

- `204` - no withdrawals found.
- `401` — user not authorized.
- `500` — internal server error.

### Interaction with Loyalty Points Calculation System

For interaction with the system, one handler is available:

- `GET /api/orders/{number}` — retrieving information about loyalty points calculation.

Request format:

```
GET /api/orders/{number} HTTP/1.1
Content-Length: 0
```

Possible response codes:

- `200` — successful request processing.

  Response format:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "order": "<number>",
        "status": "PROCESSED",
        "accrual": 500
    }
    ```

  Response object fields:

    - `order` — order number;
    - `status` — calculation status:

        - `REGISTERED` — order registered, but accrual not calculated;
        - `INVALID` — order not accepted for calculation, and reward will not be credited;
        - `PROCESSING` — accrual calculation in progress;
        - `PROCESSED` — accrual calculation completed;

    - `accrual` — calculated points to be credited, if no accrual — field is absent in response.

- `204` - order not registered in calculation system.

- `429` — service request limit exceeded.

  Response format:

    ```
    429 Too Many Requests HTTP/1.1
    Content-Type: text/plain
    Retry-After: 60
    
    No more than N requests per minute allowed
    ```

- `500` — internal server error.

An order can be taken into calculation at any time after its completion. System calculation execution time is not regulated. Statuses `INVALID` and `PROCESSED` are final.

The total number of accrual information requests is not limited.

### Loyalty System Service Configuration

The service must support configuration using the following methods:

- service startup address and port: OS environment variable `RUN_ADDRESS` or flag `-a`
- database connection address: OS environment variable `DATABASE_URI` or flag `-d`
- accrual calculation system address: OS environment variable `ACCRUAL_SYSTEM_ADDRESS` or flag `-r`
