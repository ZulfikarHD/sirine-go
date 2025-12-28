# doc-02-split-feature

I have a feature documentation file that's too large (>300 lines). Help me split it properly.

Please help me:
1. Identify what content belongs in the main feature doc (≤300 lines):
   - Overview
   - User stories
   - Business rules
   - Technical implementation summary
   - Edge cases

2. Identify what should be extracted to separate files:
   - User journey diagrams → docs/guides/{feature}-user-journeys.md
   - Full test cases → docs/testing/{CODE}-test-plan.md
   - Complete API routes → docs/api/{resource}.md

3. Create the split files with proper cross-references

Output:
- Lean feature doc (≤300 lines)
- Separate user journeys file
- Separate test plan file
- Separate API doc file (if applicable)
- Cross-reference links between all files


Current file: 