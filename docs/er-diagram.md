# ER図

`internal/shared/model` 配下のモデルのER図です。

```mermaid
erDiagram
    User {
        string ID PK
        string Email
        string Grade "nullable"
        string Course "nullable"
        string Class "nullable"
    }

    Faculty {
        uuid ID PK
        string Name
        string Email
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Room {
        uuid ID PK
        string Name
        string Floor
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Syllabus {
        string ID PK
        string Name
        string EnName
        string Grades
        int Credit
        string FacultyNames
        string PracticalHomeFacultyCategory
        string MultiplePersonTeachingForm
        string TeachingForm
        string Summary
        string LearningOutcomes
        string Assignments
        string EvaluationMethod
        string Textbooks
        string ReferenceBooks
        string Prerequisites
        string PreLearning
        string PostLearning
        string Notes
        string Keywords
        string TargetCourses
        string TargetAreas
        string Classifications
        string TeachingLanguage
        string ContentsAndSchedule
        string TeachingAndExamForm
        string DsopSubject
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Subject {
        uuid ID PK
        string Name
        int Year
        string Semester
        int Credit
        string Classification
        string CulturalSubjectCategory
        string SyllabusID FK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    SubjectFaculty {
        uuid ID PK
        uuid SubjectID FK
        uuid FacultyID FK
        bool IsPrimary
    }

    SubjectEligibleAttribute {
        uuid ID PK
        uuid SubjectID FK
        string Grade
        string Class "nullable"
    }

    SubjectRequirement {
        uuid ID PK
        uuid SubjectID FK
        string Course
        string RequirementType
    }

    TimetableItem {
        uuid ID PK
        uuid SubjectID FK
        string DayOfWeek "nullable"
        string Period "nullable"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    TimetableItemRoom {
        uuid ID PK
        uuid TimetableItemID FK
        uuid RoomID FK
    }

    CourseRegistration {
        string UserID PK "FK"
        uuid SubjectID PK "FK"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    FacultyRoom {
        uuid FacultyID PK "FK"
        uuid RoomID PK "FK"
        int Year PK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    RoomChange {
        uuid ID PK
        uuid SubjectID FK
        date Date
        string Period
        uuid OriginalRoomID FK
        uuid NewRoomID FK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    CancelledClass {
        uuid ID PK
        uuid SubjectID FK
        date Date
        string Period
        string Comment "nullable"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    MakeupClass {
        uuid ID PK
        uuid SubjectID FK
        date Date
        string Period
        string Comment "nullable"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Announcement {
        uuid ID PK
        string Title
        string URL
        timestamp AvailableFrom
        timestamp AvailableUntil "nullable"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Notification {
        string ID PK
        string Title
        string Body
        string ImageURL "nullable"
        string AnalyticsLabel "nullable"
        int APNsBadge "nullable"
        string APNsSound "nullable"
        bool APNsContentAvailable "nullable"
        string AndroidChannelID "nullable"
        string AndroidPriority "nullable"
        int AndroidTTLSeconds "nullable"
        string WebpushLink "nullable"
        string URL "nullable"
        timestamp NotifyAfter
        timestamp NotifyBefore
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    NotificationTargetUser {
        uuid NotificationID PK "FK"
        string UserID PK "FK"
        timestamp NotifiedAt "nullable"
    }

    FCMToken {
        string Token PK
        string UserID FK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    RoomReservations{
        uuid Id PK
        uuid RoomId FK
        string Title
        timestamp StartTime
        timestamp EndTime
    }

    routes["routes.txt — 路線"] {
        string id PK
        string route_id
        string route_short_name
    }

    calendar["calendar.txt — 運行曜日"] {
        string id PK
        string service_id
        int monday
        int tuesday
        int wednesday
        int thursday
        int friday
        int saturday
        int sunday
        string start_date
        string end_date
    }

    calendar_dates["calendar_dates.txt — 運行日の例外"] {
        string id PK
        string service_id FK
        string date
        int exception_type
    }

    trips["trips.txt — 便"] {
        string id PK
        string trip_id
        string route_id FK
        string service_id FK
        int direction_id
    }

    stops["stops.txt — 停留所"] {
        string id PK
        string stop_id
        string stop_name
    }

    stop_times["stop_times.txt — 停車時刻"] {
        string id PK
        string trip_id FK
        string arrival_time
        string departure_time
        string stop_id FK
        int stop_sequence
    }

    fare_rules["fare_rules.txt — 運賃ルール"] {
        string id PK
        string fare_id FK
        string route_id FK
        string origin_id FK
        string destination_id FK
        float price
    }

    Subject ||--|| Syllabus : "has"
    Subject ||--o{ SubjectFaculty : "has"
    Subject ||--o{ SubjectEligibleAttribute : "has"
    Subject ||--o{ SubjectRequirement : "has"
    SubjectFaculty }o--|| Faculty : "references"
    TimetableItem }o--|| Subject : "belongs to"
    TimetableItem ||--o{ TimetableItemRoom : "has"
    TimetableItemRoom }o--|| Room : "references"
    CourseRegistration }o--|| User : "belongs to"
    CourseRegistration }o--|| Subject : "belongs to"
    FacultyRoom }o--|| Faculty : "belongs to"
    FacultyRoom }o--|| Room : "belongs to"
    RoomChange }o--|| Subject : "belongs to"
    RoomChange }o--|| Room : "original room"
    RoomChange }o--|| Room : "new room"
    CancelledClass }o--|| Subject : "belongs to"
    MakeupClass }o--|| Subject : "belongs to"
    NotificationTargetUser }o--|| Notification : "belongs to"
    NotificationTargetUser }o--|| User : "belongs to"
    FCMToken }o--|| User : "belongs to"
    ReservedRooms }o--|| Room : "belongs to"

    routes ||--o{ trips : route_id
    calendar ||--o{ trips : service_id
    calendar ||--o{ calendar_dates : service_id
    trips ||--o{ stop_times : trip_id
    stops ||--o{ stop_times : stop_id
    routes ||--o{ fare_rules : route_id
    stops ||--o{ fare_rules : origin_id
    stops ||--o{ fare_rules : destination_id
```
