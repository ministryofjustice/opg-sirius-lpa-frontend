const initCalendar = () => {
  const calendarContainers = document.querySelectorAll(
    '[data-module="calendar"]',
  );

  if (calendarContainers.length === 0) {
    return;
  }

  // Initialize each calendar container
  calendarContainers.forEach((container, index) => {
    let bankHolidays = {};
    try {
      const bankHolidaysData = container.getAttribute("data-bank-holidays");
      if (bankHolidaysData) {
        bankHolidays = JSON.parse(bankHolidaysData);
      }
    } catch (e) {
      console.error("Failed to parse bank holidays:", e);
    }

    const now = new Date();
    const startMonth = now.getMonth() - index; // First calendar shows current month, second shows prev, third shows next
    const startYear = now.getFullYear();

    renderCalendars(container, startMonth, startYear, bankHolidays);
  });
};

const renderCalendars = (container, startMonth, startYear, bankHolidays) => {
  const calendars = [];
  const offsets = [-1, 0, 1];
  const calendarMonths = []; // Track months for each calendar

  offsets.forEach((offset, index) => {
    let month = startMonth + offset;
    let year = startYear;

    // Handle year wrapping
    if (month < 0) {
      month = 11;
      year--;
    } else if (month > 11) {
      month = 0;
      year++;
    }

    calendarMonths.push({ month, year, index });
    calendars.push(renderMonth(month, year, bankHolidays, index));
  });

  container.innerHTML = "";

  const calendarGrid = document.createElement("div");
  calendarGrid.className = "panel-calendar";
  calendarGrid.innerHTML = calendars.join("");

  container.appendChild(calendarGrid);

  calendarMonths.forEach(({ month, year, index }) => {
    const prevButton = container.querySelector(`.prev-month-${index}`);
    const nextButton = container.querySelector(`.next-month-${index}`);

    if (prevButton) {
      prevButton.addEventListener("click", () => {
        const newMonth = month - 1;
        let newYear = year;
        if (newMonth < 0) {
          newYear = year - 1;
        }
        updateCalendar(
          container,
          index,
          newMonth < 0 ? 11 : newMonth,
          newYear,
          bankHolidays,
        );
      });
    }

    if (nextButton) {
      nextButton.addEventListener("click", () => {
        const newMonth = month + 1;
        let newYear = year;
        if (newMonth > 11) {
          newYear = year + 1;
        }
        updateCalendar(
          container,
          index,
          newMonth > 11 ? 0 : newMonth,
          newYear,
          bankHolidays,
        );
      });
    }
  });
};

const updateCalendar = (
  container,
  calendarIndex,
  month,
  year,
  bankHolidays,
) => {
  const opgCalendars = container.querySelectorAll("opg-calendar");
  if (opgCalendars[calendarIndex]) {
    // Re-render with updated state
    const newHtml = renderMonth(month, year, bankHolidays, calendarIndex);
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = newHtml;
    const newOgpCalendar = tempDiv.firstChild;

    // Replace the old calendar with the new one
    opgCalendars[calendarIndex].parentNode.replaceChild(
      newOgpCalendar,
      opgCalendars[calendarIndex],
    );

    const prevButton = newOgpCalendar.querySelector(
      `.prev-month-${calendarIndex}`,
    );
    const nextButton = newOgpCalendar.querySelector(
      `.next-month-${calendarIndex}`,
    );

    if (prevButton) {
      prevButton.addEventListener("click", () => {
        const newMonth = month - 1;
        let newYear = year;
        if (newMonth < 0) {
          newYear = year - 1;
        }
        updateCalendar(
          container,
          calendarIndex,
          newMonth < 0 ? 11 : newMonth,
          newYear,
          bankHolidays,
        );
      });
    }

    if (nextButton) {
      nextButton.addEventListener("click", () => {
        const newMonth = month + 1;
        let newYear = year;
        if (newMonth > 11) {
          newYear = year + 1;
        }
        updateCalendar(
          container,
          calendarIndex,
          newMonth > 11 ? 0 : newMonth,
          newYear,
          bankHolidays,
        );
      });
    }
  }
};

const renderMonth = (month, year, bankHolidays, index) => {
  const monthNames = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];

  const dayNames = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

  const firstDay = new Date(year, month, 1).getDay();
  const adjustedFirstDay = (firstDay - 1 + 7) % 7;
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const today = new Date();
  const isCurrentMonth =
    month === today.getMonth() && year === today.getFullYear();

  const isBankHoliday = (day) => {
    if (!bankHolidays || !bankHolidays[year]) return false;

    const monthStr = String(month + 1).padStart(2, "0");
    const dayStr = String(day).padStart(2, "0");
    const dateStr = `${year}-${monthStr}-${dayStr}`;

    for (const holidayName in bankHolidays[year]) {
      const holidayDate = bankHolidays[year][holidayName].split("T")[0];
      if (holidayDate === dateStr) {
        return true;
      }
    }
    return false;
  };

  let html = `<opg-calendar>
        <div class="calendar">
         <div class="current-month">
            <div class="move-month prev-month-${index}">←</div>
            <span>${monthNames[month]} ${year}</span>
            <div class="move-month next-month-${index}">→</div>  
         </div>
    `;

  html += '<div class="week">';
  dayNames.forEach((day) => {
    html += `<div class="weekday">${day}</div>`;
  });
  html += "</div>";

  let dayCount = 0;
  html += '<div class="week">';
  for (let i = 0; i < adjustedFirstDay; i++) {
    html +=
      '<div class="day default disabled"><div class="day-number"></div><div class="event-title"></div></div>';
    dayCount++;

    if (dayCount % 7 === 0) {
      html += '</div><div class="week">';
    }
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const isToday = isCurrentMonth && day === today.getDate();
    const isLastDay = day === daysInMonth;
    const isBankHolidayDay = isBankHoliday(day);
    const dayClass = isBankHolidayDay
      ? "default disabled"
      : isToday
        ? "default"
        : "";
    const bankHolidayAttr = isBankHolidayDay ? 'data-bank-holiday="true"' : "";
    html += `<div class="day ${dayClass}" ${isLastDay ? 'data-last-day="true"' : ""} ${bankHolidayAttr}><div class="day-number">${day}</div><div class="event-title"></div></div>`;
    dayCount++;

    if (dayCount % 7 === 0 && dayCount < 35) {
      html += '</div><div class="week">';
    }
  }

  const trailingDays = 35 - dayCount;
  for (let i = 0; i < trailingDays; i++) {
    html +=
      '<div class="day default disabled"><div class="day-number"></div><div class="event-title"></div></div>';
    dayCount++;

    if (dayCount % 7 === 0 && dayCount < 35) {
      html += '</div><div class="week">';
    }
  }

  html += "</div>";

  html += "</div></opg-calendar>";

  return html;
};

export default initCalendar;
