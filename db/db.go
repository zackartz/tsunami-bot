package db

var Client *PrismaClient

func Init() error {
	Client = NewClient()
	if err := Client.Prisma.Connect(); err != nil {
		return err
	}

	return nil
}
