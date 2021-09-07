package db

var Client *PrismaClient

func Init() error {
	Client = NewClient()
	if err := Client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := Client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	return nil
}
