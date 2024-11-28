# pdf_processor.py
import json
import sys
import PyPDF2

def process_pdf(options_json):
    """
    Process a PDF file according to the provided options.

    Args:
        options_json (str): JSON string containing processing options

    Returns:
        bool: True if processing was successful
    """
    try:
        # Parse options
        options = json.loads(options_json)
        input_path = options['input_path']
        output_path = options['output_path']
        metadata = options['metadata']
        # Open the PDF
        with open(input_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)
            writer = PyPDF2.PdfWriter()

            # Copy all pages
            for page in reader.pages:
                writer.add_page(page)

            # Update metadata
            writer.update_page_form_field_values(
                writer.pages[0], metadata
            )

            # Save the modified PDF
            with open(output_path, 'wb') as output_file:
                writer.write(output_file)

        return True

    except Exception as e:
        print(f"Error processing PDF: {str(e)}")
        return False

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python pdf_processor.py '<json_string>'")
        sys.exit(1)

    process_pdf(sys.argv[1])