# developement-strategy-build

I need you to analyze a feature development file and create a comprehensive cross-frontend implementation strategy.
---

ANALYSIS REQUIREMENTS:

PHASE 1: FEATURE UNDERSTANDING
1. Read the feature specification thoroughly
2. Identify:
   - What data is being created/managed?
   - Who creates this data? (Owner)
   - Who needs to see/use this data? (Consumer)
   - What's the primary user goal?

PHASE 2: CROSS-FRONTEND IMPACT MAPPING
For EACH feature in the file, create a table:

| Feature Name | Owner (Who Creates) | Consumer (Who Views) | Data Flow |
|--------------|----------------------|---------------------|-----------|
| [Feature]    | [Frontend + Action]  | [Frontend + View]   | Create→Store→Display |

PHASE 3: MISSING IMPLEMENTATION DETECTION
For each feature, check if the specification covers:

**Owner Side (Data Creation):**
- [ ] UI form/interface for creating data
- [ ] Validation rules
- [ ] Edit/Update capability
- [ ] Delete/Archive capability
- [ ] Preview before publishing
- [ ] Bulk operations (if applicable)

**Consumer Side (Data Display):**
- [ ] Where users will SEE this data (specific page/component)
- [ ] How users will FIND this data (navigation/search)
- [ ] What users can DO with this data (interactions)
- [ ] Mobile/responsive version
- [ ] Empty states (when no data exists)
- [ ] Loading states

**Integration Points:**
- [ ] API endpoints needed
- [ ] Database schema changes
- [ ] State management updates
- [ ] Navigation menu updates
- [ ] Search/filter additions
- [ ] Notification/real-time updates

PHASE 4: GAP ANALYSIS
Identify and flag ANY features where:
⚠️ Owner can create data BUT Consumer has no way to view it
⚠️ Consumer expects to see data BUT Owner has no way to create it
⚠️ Feature exists in backend BUT no frontend UI specified
⚠️ Feature exists in one frontend BUT missing in others where it's needed

PHASE 5: IMPLEMENTATION SEQUENCING
For each feature, determine:
1. **What must be built FIRST** (dependencies)
2. **What can be built in PARALLEL** (independent work)
3. **What should be built LAST** (dependent on other work)

Create a priority matrix:
- P0 (Critical): Feature is unusable without this
- P1 (Important): Feature is incomplete without this
- P2 (Enhancement): Nice to have, can ship later

PHASE 6: DETAILED RECOMMENDATIONS
For each frontend that needs updates, specify:

**New Pages/Routes Needed:**
- [ ] Page: [Name] - Purpose: [Why] - Priority: [P0/P1/P2]

**Updates to Existing Pages:**
- [ ] Page: [Name] - Add: [What] - Location: [Where] - Priority: [P0/P1/P2]

**Navigation/Menu Changes:**
- [ ] Frontend: [Which] - Add Menu Item: [Name] - Parent: [Where] - Priority: [P0/P1/P2]

**Component Library Additions:**
- [ ] Component: [Name] - Used By: [Pages] - Priority: [P0/P1/P2]

PHASE 7: EXAMPLE USER JOURNEYS
For the TOP 3 most important features, write out:

**Owner Journey:**
1. User navigates to: [Specific path]
2. User clicks: [Button/Link]
3. User fills: [Form fields]
4. System does: [Backend action]
5. User sees: [Feedback/Result]

**Consumer Journey:**
1. User navigates to: [Specific path]
2. User sees: [What's displayed]
3. User can: [Available actions]
4. When user clicks: [What happens]
5. User achieves: [Goal completed]

---

CRITICAL REMINDERS:
- Think holistically: data creation WITHOUT display = broken feature
- Check mobile: if desktop has it, mobile probably needs it too
- Verify navigation: users must be able to FIND what owners create
- Consider empty states: what if no data exists yet?
- Think about filters: if searchable data exists, add search/filter UI
- Don't assume: if it's not explicitly mentioned, flag it as missing

Now analyze the file and provide the complete analysis.