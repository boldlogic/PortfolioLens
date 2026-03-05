
SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[fx_rates](
	[date] [date] NOT NULL,
	[quote_iso_code] [bigint] NOT NULL,
	[base_iso_code] [bigint] NOT NULL,
	[nominal] [bigint] NULL,
	[quote_for_nominal] [float] NULL,
	[quote_per_unit] [float] NULL,
	[base_per_quote_unit] [float] NULL,
	[created_at] [datetimeoffset](7) NULL,
	[updated_at] [datetimeoffset](7) NULL,
PRIMARY KEY CLUSTERED 
(
	[date] ASC,
	[quote_iso_code] ASC,
	[base_iso_code] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO


